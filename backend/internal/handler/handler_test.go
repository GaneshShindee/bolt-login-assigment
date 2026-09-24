package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/ganeshshinde/boltapp/backend/internal/auth"
	"github.com/ganeshshinde/boltapp/backend/internal/biz"
	"github.com/ganeshshinde/boltapp/backend/internal/repotest"
	"github.com/ganeshshinde/boltapp/backend/internal/service"
)

type testAPI struct {
	t         *testing.T
	h         http.Handler
	users     *repotest.Users
	checkouts *repotest.Checkouts
}

// newTestAPI wires the real handler → service → biz layers on in-memory repositories.
func newTestAPI(t *testing.T) *testAPI {
	users, checkouts := &repotest.Users{}, &repotest.Checkouts{}
	userBiz := biz.NewUserBiz(users, auth.NewTokenSigner("test-secret", time.Hour), auth.NewLoginLimiter(5, 15*time.Minute))
	srv := NewServer(
		service.NewUserService(userBiz),
		service.NewCheckoutService(userBiz, biz.NewCheckoutBiz(checkouts)),
		func(context.Context) error { return nil },
	)
	return &testAPI{t: t, h: srv.Routes([]string{"http://localhost:5173"}), users: users, checkouts: checkouts}
}

// do sends a JSON request and decodes the JSON response into a map.
func (a *testAPI) do(method, path, body, token string) (int, map[string]any) {
	a.t.Helper()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	rec := httptest.NewRecorder()
	a.h.ServeHTTP(rec, req)
	var out map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		a.t.Fatalf("%s %s: invalid JSON response %q", method, path, rec.Body.String())
	}
	return rec.Code, out
}

// TestLoginFlow walks the whole happy path and the key failure cases:
// register → recognize → wrong code → right code → /me → checkout linked to the user.
func TestLoginFlow(t *testing.T) {
	api := newTestAPI(t)

	// Register (email is normalized).
	status, body := api.do("POST", "/api/register",
		`{"email":"  Priya@Example.com ","first_name":"Priya","last_name":"Sharma"}`, "")
	if status != http.StatusCreated {
		t.Fatalf("register: status %d, body %v", status, body)
	}
	code, _ := body["code"].(string)
	if len(code) != 6 {
		t.Fatalf("register: expected a 6-digit code, got %q", code)
	}
	if hash := api.users.All[0].LoginCodeHash; hash == code || !auth.CheckCode(hash, code) {
		t.Fatal("register: code must be stored as a bcrypt hash")
	}
	if user := body["user"].(map[string]any); len(user) != 4 {
		t.Fatalf("register: user should only expose id, email and names, got %v", user)
	}

	// Duplicate registration.
	status, _ = api.do("POST", "/api/register", `{"email":"priya@example.com","first_name":"P","last_name":"S"}`, "")
	if status != http.StatusConflict {
		t.Fatalf("duplicate register: status %d, want 409", status)
	}

	// Recognize: known vs unknown, and no personal data leaked.
	_, body = api.do("POST", "/api/recognize", `{"email":"PRIYA@example.com"}`, "")
	if body["recognized"] != true || len(body) != 1 {
		t.Fatalf("recognize known: got %v, want only {recognized:true}", body)
	}
	_, body = api.do("POST", "/api/recognize", `{"email":"nobody@example.com"}`, "")
	if body["recognized"] != false {
		t.Fatalf("recognize unknown: got %v", body)
	}

	// Wrong code → 401 with an error message.
	wrong := "000000"
	if code == wrong {
		wrong = "111111"
	}
	status, body = api.do("POST", "/api/login", `{"email":"priya@example.com","code":"`+wrong+`"}`, "")
	if status != http.StatusUnauthorized || body["error"] == nil {
		t.Fatalf("wrong code: status %d, body %v", status, body)
	}

	// Right code → token + user.
	status, body = api.do("POST", "/api/login", `{"email":"priya@example.com","code":"`+code+`"}`, "")
	if status != http.StatusOK {
		t.Fatalf("login: status %d, body %v", status, body)
	}
	token, _ := body["token"].(string)
	if token == "" {
		t.Fatal("login: missing token")
	}
	if name := body["user"].(map[string]any)["first_name"]; name != "Priya" {
		t.Fatalf("login: first_name = %v", name)
	}

	// /me with and without the token.
	if status, _ := api.do("GET", "/api/me", "", token); status != http.StatusOK {
		t.Fatalf("me with token: status %d", status)
	}
	if status, _ := api.do("GET", "/api/me", "", ""); status != http.StatusUnauthorized {
		t.Fatalf("me without token: status %d", status)
	}

	// Checkout while logged in is linked to the user; as a guest it is not.
	checkout := `{"email":"priya@example.com","phone":"9876543210","address":{"label":"Home","line1":"221B Baker St","line2":"","city":"Pune","state":"Maharashtra","pincode":"411001"}}`
	if status, body := api.do("POST", "/api/checkout", checkout, token); status != http.StatusCreated {
		t.Fatalf("checkout: status %d, body %v", status, body)
	}
	if status, _ := api.do("POST", "/api/checkout", checkout, ""); status != http.StatusCreated {
		t.Fatalf("guest checkout: status %d", status)
	}
	if c := api.checkouts.All[0]; c.UserID == nil || *c.UserID != 1 {
		t.Fatalf("logged-in checkout should have user_id 1, got %v", c.UserID)
	}
	if c := api.checkouts.All[1]; c.UserID != nil {
		t.Fatalf("guest checkout should have no user_id, got %v", *c.UserID)
	}

	// Saved details: only the logged-in order's phone + address, and only with a token.
	status, body = api.do("GET", "/api/me/saved-details", "", token)
	saved, _ := body["saved"].([]any)
	if status != http.StatusOK || len(saved) != 1 {
		t.Fatalf("saved details: status %d, body %v", status, body)
	}
	if city := saved[0].(map[string]any)["address"].(map[string]any)["city"]; city != "Pune" {
		t.Fatalf("saved details: city = %v", city)
	}
	if status, _ := api.do("GET", "/api/me/saved-details", "", ""); status != http.StatusUnauthorized {
		t.Fatalf("saved details without token: status %d, want 401", status)
	}
}

func TestLogin_RateLimited(t *testing.T) {
	api := newTestAPI(t)
	_, body := api.do("POST", "/api/register", `{"email":"a@example.com","first_name":"A","last_name":"B"}`, "")
	code := body["code"].(string)
	wrong := "000000"
	if code == wrong {
		wrong = "111111"
	}

	for i := range 5 {
		if status, _ := api.do("POST", "/api/login", `{"email":"a@example.com","code":"`+wrong+`"}`, ""); status != http.StatusUnauthorized {
			t.Fatalf("attempt %d: status %d, want 401", i+1, status)
		}
	}
	// Even the correct code is refused once the limit is hit.
	if status, _ := api.do("POST", "/api/login", `{"email":"a@example.com","code":"`+code+`"}`, ""); status != http.StatusTooManyRequests {
		t.Fatalf("after 5 failures: status %d, want 429", status)
	}
}

func TestValidationErrors(t *testing.T) {
	api := newTestAPI(t)
	tests := []struct{ name, path, body string }{
		{"register bad email", "/api/register", `{"email":"a@b","first_name":"A","last_name":"B"}`},
		{"register missing name", "/api/register", `{"email":"a@b.com","first_name":" ","last_name":"B"}`},
		{"login short code", "/api/login", `{"email":"a@b.com","code":"123"}`},
		{"checkout bad phone", "/api/checkout", `{"email":"a@b.com","phone":"abc","address":{"label":"Home","line1":"221B Baker St","line2":"","city":"Pune","state":"Maharashtra","pincode":"411001"}}`},
		{"checkout no address", "/api/checkout", `{"email":"a@b.com","phone":"9876543210"}`},
		{"checkout bad pincode", "/api/checkout", `{"email":"a@b.com","phone":"9876543210","address":{"label":"Home","line1":"x","city":"y","state":"z","pincode":"12"}}`},
		{"unknown field", "/api/recognize", `{"email":"a@b.com","extra":1}`},
		{"not json", "/api/recognize", `not json`},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			status, body := api.do("POST", tc.path, tc.body, "")
			if status != http.StatusBadRequest || body["error"] == nil {
				t.Fatalf("status %d, body %v; want 400 with error", status, body)
			}
		})
	}
}

func TestCORS(t *testing.T) {
	api := newTestAPI(t)
	req := httptest.NewRequest("OPTIONS", "/api/login", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	rec := httptest.NewRecorder()
	api.h.ServeHTTP(rec, req)
	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "http://localhost:5173" {
		t.Fatalf("allowed origin: got %q", got)
	}

	req = httptest.NewRequest("OPTIONS", "/api/login", nil)
	req.Header.Set("Origin", "https://evil.example")
	rec = httptest.NewRecorder()
	api.h.ServeHTTP(rec, req)
	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Fatalf("disallowed origin should get no CORS header, got %q", got)
	}
}
