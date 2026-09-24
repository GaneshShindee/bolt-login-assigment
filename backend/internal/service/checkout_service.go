package service

import (
	"context"

	"github.com/ganeshshinde/boltapp/backend/internal/biz"
	"github.com/ganeshshinde/boltapp/backend/internal/entity"
)

// CheckoutService combines two biz components: identifying the (optional) logged-in user,
// then recording the order.
type CheckoutService struct {
	users     *biz.UserBiz
	checkouts *biz.CheckoutBiz
}

func NewCheckoutService(users *biz.UserBiz, checkouts *biz.CheckoutBiz) *CheckoutService {
	return &CheckoutService{users: users, checkouts: checkouts}
}

// Checkout links the order to the user for a valid token; with no token or an expired one it
// is saved as a guest checkout.
func (s *CheckoutService) Checkout(ctx context.Context, token string, req CheckoutRequest) (CheckoutResponse, error) {
	var userID *int
	if token != "" {
		if u, err := s.users.Authenticate(ctx, token); err == nil {
			userID = &u.ID
		}
	}
	c, err := s.checkouts.Place(ctx, entity.Checkout{
		UserID:  userID,
		Email:   req.Email,
		Phone:   req.Phone,
		Address: req.Address.toEntity(),
	})
	if err != nil {
		return CheckoutResponse{}, err
	}
	return toCheckoutResponse(c), nil
}

func (s *CheckoutService) SavedDetails(ctx context.Context, token string) (SavedDetailsResponse, error) {
	u, err := s.users.Authenticate(ctx, token)
	if err != nil {
		return SavedDetailsResponse{}, err
	}
	details, err := s.checkouts.SavedDetails(ctx, u.ID)
	if err != nil {
		return SavedDetailsResponse{}, err
	}
	resp := SavedDetailsResponse{Saved: make([]SavedDetailsDTO, 0, len(details))}
	for _, d := range details {
		resp.Saved = append(resp.Saved, SavedDetailsDTO{Phone: d.Phone, Address: toAddressDTO(d.Address), LastUsedAt: d.LastUsedAt})
	}
	return resp, nil
}
