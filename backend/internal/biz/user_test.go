package biz_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ganeshshinde/boltapp/backend/internal/auth"
	"github.com/ganeshshinde/boltapp/backend/internal/biz"
	"github.com/ganeshshinde/boltapp/backend/internal/repotest"
)

func newUserBiz() (*biz.UserBiz, *repotest.Users) {
	repo := &repotest.Users{}
	return biz.NewUserBiz(repo, auth.NewTokenSigner("test", time.Hour), auth.NewLoginLimiter(3, time.Minute)), repo
}

func wrongCode(code string) string {
	if code == "000000" {
		return "111111"
	}
	return "000000"
}

func TestRegister_NormalizesAndHashes(t *testing.T) {
	b, repo := newUserBiz()
	u, code, err := b.Register(context.Background(), "  Priya@Example.COM ", " Priya ", " Sharma ")
	if err != nil {
		t.Fatalf("Register: %v", err)
	}
	if u.Email != "priya@example.com" || u.FirstName != "Priya" || u.LastName != "Sharma" {
		t.Fatalf("fields not normalized: %+v", u)
	}
	stored := repo.All[0]
	if stored.LoginCodeHash == code || !auth.CheckCode(stored.LoginCodeHash, code) {
		t.Fatal("only a bcrypt hash of the code should be stored")
	}
}

func TestRegister_Errors(t *testing.T) {
	b, _ := newUserBiz()
	ctx := context.Background()
	if _, _, err := b.Register(ctx, "a@example.com", "A", "B"); err != nil {
		t.Fatal(err)
	}
	if _, _, err := b.Register(ctx, "A@EXAMPLE.com", "A", "B"); !errors.Is(err, biz.ErrEmailTaken) {
		t.Fatalf("duplicate (case-insensitive): got %v, want ErrEmailTaken", err)
	}

	var ve *biz.ValidationError
	for _, tc := range []struct{ email, first, last string }{
		{"not-an-email", "A", "B"},
		{"b@example.com", "   ", "B"},
		{"b@example.com", "A", ""},
	} {
		if _, _, err := b.Register(ctx, tc.email, tc.first, tc.last); !errors.As(err, &ve) {
			t.Errorf("Register(%q, %q, %q): got %v, want ValidationError", tc.email, tc.first, tc.last, err)
		}
	}
}

func TestRecognize(t *testing.T) {
	b, _ := newUserBiz()
	ctx := context.Background()
	_, _, _ = b.Register(ctx, "a@example.com", "A", "B")

	if ok, err := b.Recognize(ctx, " A@Example.com"); err != nil || !ok {
		t.Fatalf("registered email: got %v, %v", ok, err)
	}
	if ok, err := b.Recognize(ctx, "nobody@example.com"); err != nil || ok {
		t.Fatalf("unknown email: got %v, %v", ok, err)
	}
	var ve *biz.ValidationError
	if _, err := b.Recognize(ctx, "a@b"); !errors.As(err, &ve) {
		t.Fatalf("incomplete email: got %v, want ValidationError", err)
	}
}

func TestLogin(t *testing.T) {
	b, _ := newUserBiz()
	ctx := context.Background()
	registered, code, _ := b.Register(ctx, "a@example.com", "A", "B")

	if _, _, err := b.Login(ctx, "a@example.com", wrongCode(code)); !errors.Is(err, biz.ErrInvalidCode) {
		t.Fatalf("wrong code: got %v, want ErrInvalidCode", err)
	}
	if _, _, err := b.Login(ctx, "nobody@example.com", code); !errors.Is(err, biz.ErrInvalidCode) {
		t.Fatalf("unknown email: got %v, want ErrInvalidCode", err)
	}

	u, token, err := b.Login(ctx, "A@example.com", code)
	if err != nil || u.ID != registered.ID || token == "" {
		t.Fatalf("right code: got %+v, %q, %v", u, token, err)
	}
	if got, err := b.Authenticate(ctx, token); err != nil || got.ID != registered.ID {
		t.Fatalf("Authenticate: got %+v, %v", got, err)
	}
	if _, err := b.Authenticate(ctx, "garbage"); !errors.Is(err, biz.ErrUnauthenticated) {
		t.Fatalf("bad token: got %v, want ErrUnauthenticated", err)
	}
}

func TestLogin_RateLimited(t *testing.T) {
	b, _ := newUserBiz() // limit: 3 failures per minute
	ctx := context.Background()
	_, code, _ := b.Register(ctx, "a@example.com", "A", "B")

	for range 3 {
		_, _, _ = b.Login(ctx, "a@example.com", wrongCode(code))
	}
	if _, _, err := b.Login(ctx, "a@example.com", code); !errors.Is(err, biz.ErrTooManyAttempts) {
		t.Fatalf("after 3 failures even the right code is refused: got %v", err)
	}
}
