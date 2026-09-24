package biz

import (
	"context"
	"strings"

	"github.com/ganeshshinde/boltapp/backend/internal/entity"
)

const savedDetailsToShow = 3

type CheckoutBiz struct {
	repo CheckoutRepo
}

func NewCheckoutBiz(repo CheckoutRepo) *CheckoutBiz {
	return &CheckoutBiz{repo: repo}
}

func (b *CheckoutBiz) Place(ctx context.Context, c entity.Checkout) (entity.Checkout, error) {
	c.Email = normalizeEmail(c.Email)
	c.Phone = strings.TrimSpace(c.Phone)
	c.Address = normalizeAddress(c.Address)
	if err := validateCheckout(c); err != nil {
		return entity.Checkout{}, err
	}
	return b.repo.Create(ctx, c)
}

func (b *CheckoutBiz) SavedDetails(ctx context.Context, userID int) ([]entity.SavedDetails, error) {
	return b.repo.RecentDetails(ctx, userID, savedDetailsToShow)
}

func normalizeAddress(a entity.Address) entity.Address {
	return entity.Address{
		Label:   strings.TrimSpace(a.Label),
		Line1:   strings.TrimSpace(a.Line1),
		Line2:   strings.TrimSpace(a.Line2),
		City:    strings.TrimSpace(a.City),
		State:   strings.TrimSpace(a.State),
		Pincode: strings.TrimSpace(a.Pincode),
	}
}

func validateCheckout(c entity.Checkout) error {
	switch {
	case !validEmail(c.Email):
		return invalid("Please enter a valid email address.")
	case !validPhone(c.Phone):
		return invalid("Please enter a valid 10-digit mobile number.")
	case !required(c.Address.Label, maxAddressLabelLen):
		return invalid("Please give this address a name, like Home or Work.")
	case !required(c.Address.Line1, maxAddressLineLen):
		return invalid("Please enter your house / flat and street.")
	case len(c.Address.Line2) > maxAddressLineLen:
		return invalid("Address line 2 is too long.")
	case !required(c.Address.City, maxCityOrStateLen):
		return invalid("Please enter your city.")
	case !required(c.Address.State, maxCityOrStateLen):
		return invalid("Please choose your state.")
	case !validPincode(c.Address.Pincode):
		return invalid("Please enter a valid 6-digit PIN code.")
	}
	return nil
}
