// Package handler is the HTTP layer: routing, middleware, JSON and error-to-status mapping.
// It holds no business rules.
package handler

import (
	"context"
	"net/http"

	"github.com/ganeshshinde/boltapp/backend/internal/service"
)

type Server struct {
	users     *service.UserService
	checkouts *service.CheckoutService
	ping      func(context.Context) error // database health check
}

func NewServer(users *service.UserService, checkouts *service.CheckoutService, ping func(context.Context) error) *Server {
	return &Server{users: users, checkouts: checkouts, ping: ping}
}

func (s *Server) Routes(allowedOrigins []string) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", s.handleHealth)

	mux.Handle("POST /api/register", withBody(http.StatusCreated, ignoreToken(s.users.Register)))
	mux.Handle("POST /api/recognize", withBody(http.StatusOK, ignoreToken(s.users.Recognize)))
	mux.Handle("POST /api/login", withBody(http.StatusOK, ignoreToken(s.users.Login)))
	mux.Handle("GET /api/me", withoutBody(s.users.Me))

	// The session token is optional here: without one it's a guest order.
	mux.Handle("POST /api/checkout", withBody(http.StatusCreated, s.checkouts.Checkout))
	mux.Handle("GET /api/me/saved-details", withoutBody(s.checkouts.SavedDetails))

	return logRequests(cors(allowedOrigins, mux))
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	if err := s.ping(r.Context()); err != nil {
		writeError(w, http.StatusServiceUnavailable, "Database unavailable.")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
