package handler

import (
	"errors"
	"log"
	"net/http"

	"github.com/ganeshshinde/boltapp/backend/internal/biz"
)

func writeBizError(w http.ResponseWriter, err error) {
	var ve *biz.ValidationError
	switch {
	case errors.As(err, &ve):
		writeError(w, http.StatusBadRequest, ve.Message)
	case errors.Is(err, biz.ErrEmailTaken):
		writeError(w, http.StatusConflict, "An account with this email already exists.")
	case errors.Is(err, biz.ErrInvalidCode):
		writeError(w, http.StatusUnauthorized, "That code doesn't match. Please try again.")
	case errors.Is(err, biz.ErrTooManyAttempts):
		writeError(w, http.StatusTooManyRequests, "Too many failed attempts. Please wait a few minutes and try again.")
	case errors.Is(err, biz.ErrUnauthenticated):
		writeError(w, http.StatusUnauthorized, "Not logged in.")
	default:
		log.Printf("internal error: %v", err)
		writeError(w, http.StatusInternalServerError, "Something went wrong. Please try again.")
	}
}
