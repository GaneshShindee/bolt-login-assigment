package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ganeshshinde/boltapp/backend/internal/auth"
	"github.com/ganeshshinde/boltapp/backend/internal/biz"
	"github.com/ganeshshinde/boltapp/backend/internal/repotest"
	"github.com/ganeshshinde/boltapp/backend/internal/service"
)

type fixture struct {
	users     *service.UserService
	checkout  *service.CheckoutService
	checkouts *repotest.Checkouts
}

func newFixture() fixture {
	checkouts := &repotest.Checkouts{}
	userBiz := biz.NewUserBiz(&repotest.Users{}, auth.NewTokenSigner("test", time.Hour), auth.NewLoginLimiter(5, time.Minute))
	return fixture{
		users:     service.NewUserService(userBiz),
		checkout:  service.NewCheckoutService(userBiz, biz.NewCheckoutBiz(checkouts)),
		checkouts: checkouts,
	}
}

// registerAndLogin returns the new user's ID and session token.
func (f fixture) registerAndLogin(t *testing.T, email string) (int, string) {
	t.Helper()
	ctx := context.Background()
	reg, err := f.users.Register(ctx, service.RegisterRequest{Email: email, FirstName: "A", LastName: "B"})
	if err != nil {
		t.Fatal(err)
	}
	login, err := f.users.Login(ctx, service.LoginRequest{Email: email, Code: reg.Code})
	if err != nil {
		t.Fatal(err)
	}
	return reg.User.ID, login.Token
}

var testAddress = service.AddressDTO{Label: "Home", Line1: "1 Main St", City: "Pune", State: "Maharashtra", Pincode: "411001"}

// TestCheckout_LinksUserFromToken checks the orchestration across the two biz components:
// a valid token links the order to the user; no token or a bad token saves it as a guest.
func TestCheckout_LinksUserFromToken(t *testing.T) {
	f := newFixture()
	userID, token := f.registerAndLogin(t, "a@example.com")

	req := service.CheckoutRequest{Email: "a@example.com", Phone: "9876543210", Address: testAddress}
	for _, tok := range []string{token, "", "not-a-token"} {
		if _, err := f.checkout.Checkout(context.Background(), tok, req); err != nil {
			t.Fatalf("Checkout(token=%q): %v", tok, err)
		}
	}

	if c := f.checkouts.All[0]; c.UserID == nil || *c.UserID != userID {
		t.Errorf("valid token: want user_id %d, got %v", userID, c.UserID)
	}
	for i, c := range f.checkouts.All[1:] {
		if c.UserID != nil {
			t.Errorf("guest checkout %d: want no user_id, got %d", i+2, *c.UserID)
		}
	}
}

func TestSavedDetails(t *testing.T) {
	f := newFixture()
	ctx := context.Background()
	_, token := f.registerAndLogin(t, "a@example.com")

	// First visit: nothing saved yet, but an empty list (not null).
	resp, err := f.checkout.SavedDetails(ctx, token)
	if err != nil || resp.Saved == nil || len(resp.Saved) != 0 {
		t.Fatalf("before any order: got %+v, %v", resp, err)
	}

	req := service.CheckoutRequest{Email: "a@example.com", Phone: "9876543210", Address: testAddress}
	if _, err := f.checkout.Checkout(ctx, token, req); err != nil {
		t.Fatal(err)
	}
	resp, err = f.checkout.SavedDetails(ctx, token)
	if err != nil || len(resp.Saved) != 1 || resp.Saved[0].Phone != "9876543210" || resp.Saved[0].Address != testAddress {
		t.Fatalf("after one order: got %+v, %v", resp, err)
	}

	if _, err := f.checkout.SavedDetails(ctx, ""); !errors.Is(err, biz.ErrUnauthenticated) {
		t.Fatalf("no token: got %v, want ErrUnauthenticated", err)
	}
}
