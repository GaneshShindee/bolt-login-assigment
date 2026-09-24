package handler

import (
	"context"
	"net/http"
	"strings"
)

// withBody adapts a service call that takes a JSON body: decode, call, then write the result or error.
func withBody[Req, Resp any](status int, call func(ctx context.Context, token string, req Req) (Resp, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req Req
		if !readJSON(w, r, &req) {
			return
		}
		respond(w, status, func() (Resp, error) { return call(r.Context(), bearerToken(r), req) })
	}
}

func withoutBody[Resp any](call func(ctx context.Context, token string) (Resp, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		respond(w, http.StatusOK, func() (Resp, error) { return call(r.Context(), bearerToken(r)) })
	}
}

func respond[Resp any](w http.ResponseWriter, status int, call func() (Resp, error)) {
	resp, err := call()
	if err != nil {
		writeBizError(w, err)
		return
	}
	writeJSON(w, status, resp)
}

// ignoreToken lets a service call that doesn't need a session fit withBody.
func ignoreToken[Req, Resp any](call func(context.Context, Req) (Resp, error)) func(context.Context, string, Req) (Resp, error) {
	return func(ctx context.Context, _ string, req Req) (Resp, error) { return call(ctx, req) }
}

func bearerToken(r *http.Request) string {
	token, ok := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
	if !ok {
		return ""
	}
	return token
}
