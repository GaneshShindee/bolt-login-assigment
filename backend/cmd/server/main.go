package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/ganeshshinde/boltapp/backend/internal/auth"
	"github.com/ganeshshinde/boltapp/backend/internal/biz"
	"github.com/ganeshshinde/boltapp/backend/internal/handler"
	"github.com/ganeshshinde/boltapp/backend/internal/repository"
	"github.com/ganeshshinde/boltapp/backend/internal/service"
)

const (
	sessionTTL = 24 * time.Hour
	// Up to maxFailedLogins wrong codes per email within loginLockout.
	maxFailedLogins = 5
	loginLockout    = 15 * time.Minute
	shutdownTimeout = 10 * time.Second
)

func main() {
	dbURL := mustEnv("DATABASE_URL")
	secret := mustEnv("SESSION_SECRET")
	port := envOr("PORT", "8080")
	allowedOrigins := strings.Split(envOr("ALLOWED_ORIGIN", "http://localhost:5173"), ",")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	db, err := repository.Open(ctx, dbURL)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer db.Close()

	userBiz := biz.NewUserBiz(
		repository.NewUserRepo(db),
		auth.NewTokenSigner(secret, sessionTTL),
		auth.NewLoginLimiter(maxFailedLogins, loginLockout),
	)
	checkoutBiz := biz.NewCheckoutBiz(repository.NewCheckoutRepo(db))
	srv := handler.NewServer(
		service.NewUserService(userBiz),
		service.NewCheckoutService(userBiz, checkoutBiz),
		db.Ping,
	)

	httpServer := &http.Server{
		Addr:              ":" + port,
		Handler:           srv.Routes(allowedOrigins),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
	}

	// Render sends SIGTERM on redeploy: let in-flight requests finish.
	stop, cancelSignals := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancelSignals()
	shutdownDone := make(chan struct{})
	go func() {
		defer close(shutdownDone)
		<-stop.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			log.Printf("shutdown: %v", err)
		}
	}()

	log.Printf("API listening on :%s (allowed origins: %v)", port, allowedOrigins)
	if err := httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server: %v", err)
	}
	<-shutdownDone // wait for in-flight requests to drain
	log.Print("server stopped")
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("missing required env var %s", key)
	}
	return v
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
