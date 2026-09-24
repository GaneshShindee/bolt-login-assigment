package biz_test

import (
	"context"
	"errors"
	"testing"

	"github.com/ganeshshinde/boltapp/backend/internal/biz"
	"github.com/ganeshshinde/boltapp/backend/internal/entity"
	"github.com/ganeshshinde/boltapp/backend/internal/repotest"
)

func validAddress() entity.Address {
	return entity.Address{Label: "Home", Line1: "12 MG Road", Line2: "Near Metro", City: "Bengaluru", State: "Karnataka", Pincode: "560001"}
}

func TestPlace_NormalizesAndSaves(t *testing.T) {
	repo := &repotest.Checkouts{}
	b := biz.NewCheckoutBiz(repo)
	c, err := b.Place(context.Background(), entity.Checkout{
		Email: " Guest@Example.com ", Phone: " 9876543210 ",
		Address: entity.Address{Label: " Home ", Line1: "  12 MG Road ", City: " Bengaluru ", State: "Karnataka", Pincode: " 560001 "},
	})
	if err != nil {
		t.Fatalf("Place: %v", err)
	}
	want := entity.Address{Label: "Home", Line1: "12 MG Road", City: "Bengaluru", State: "Karnataka", Pincode: "560001"}
	if c.ID != 1 || c.Email != "guest@example.com" || c.Phone != "9876543210" || c.Address != want {
		t.Fatalf("not normalized or not saved: %+v", c)
	}
	if len(repo.All) != 1 {
		t.Fatalf("expected 1 saved checkout, got %d", len(repo.All))
	}
}

func TestPlace_Validation(t *testing.T) {
	repo := &repotest.Checkouts{}
	b := biz.NewCheckoutBiz(repo)
	with := func(edit func(*entity.Checkout)) entity.Checkout {
		c := entity.Checkout{Email: "a@b.com", Phone: "9876543210", Address: validAddress()}
		edit(&c)
		return c
	}
	cases := map[string]entity.Checkout{
		"bad email":         with(func(c *entity.Checkout) { c.Email = "a@b" }),
		"bad phone":         with(func(c *entity.Checkout) { c.Phone = "abc" }),
		"blank label":       with(func(c *entity.Checkout) { c.Address.Label = "  " }),
		"label too long":    with(func(c *entity.Checkout) { c.Address.Label = "My very long address nickname here" }),
		"blank line 1":      with(func(c *entity.Checkout) { c.Address.Line1 = "   " }),
		"blank city":        with(func(c *entity.Checkout) { c.Address.City = "" }),
		"blank state":       with(func(c *entity.Checkout) { c.Address.State = "" }),
		"5-digit pincode":   with(func(c *entity.Checkout) { c.Address.Pincode = "56000" }),
		"pincode leading 0": with(func(c *entity.Checkout) { c.Address.Pincode = "060001" }),
		"pincode letters":   with(func(c *entity.Checkout) { c.Address.Pincode = "56A001" }),
	}
	var ve *biz.ValidationError
	for name, c := range cases {
		if _, err := b.Place(context.Background(), c); !errors.As(err, &ve) {
			t.Errorf("%s: got %v, want ValidationError", name, err)
		}
	}
	if len(repo.All) != 0 {
		t.Fatal("invalid checkouts must not be saved")
	}
	// Line 2 is optional.
	if _, err := b.Place(context.Background(), with(func(c *entity.Checkout) { c.Address.Line2 = "" })); err != nil {
		t.Fatalf("empty line 2 should be allowed: %v", err)
	}
}

func TestSavedDetails_DistinctNewestFirst(t *testing.T) {
	repo := &repotest.Checkouts{}
	b := biz.NewCheckoutBiz(repo)
	ctx := context.Background()
	user, other := 1, 2
	home := validAddress()
	office := entity.Address{Label: "Work", Line1: "5th Floor, Tech Park", City: "Pune", State: "Maharashtra", Pincode: "411001"}

	place := func(userID *int, phone string, a entity.Address) {
		t.Helper()
		if _, err := b.Place(ctx, entity.Checkout{UserID: userID, Email: "a@b.com", Phone: phone, Address: a}); err != nil {
			t.Fatal(err)
		}
	}
	place(&user, "9876543210", home)
	place(&user, "9876543210", office)
	place(&user, "9876543210", home) // reused: should appear once, as the newest
	place(&other, "9000000000", office)
	place(nil, "9111111111", home) // guest

	got, err := b.SavedDetails(ctx, user)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 || got[0].Address != home || got[1].Address != office {
		t.Fatalf("want [home, office] for user 1, got %+v", got)
	}
	if got, _ := b.SavedDetails(ctx, 99); len(got) != 0 {
		t.Fatalf("user with no orders: want none, got %+v", got)
	}
}
