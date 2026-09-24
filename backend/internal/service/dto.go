package service

import (
	"time"

	"github.com/ganeshshinde/boltapp/backend/internal/entity"
)

type UserDTO struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func toUserDTO(u entity.User) UserDTO {
	return UserDTO{ID: u.ID, Email: u.Email, FirstName: u.FirstName, LastName: u.LastName}
}

type RegisterRequest struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type RegisterResponse struct {
	User UserDTO `json:"user"`
	Code string  `json:"code"` // shown to the user once; only its hash is stored
}

type RecognizeRequest struct {
	Email string `json:"email"`
}

type RecognizeResponse struct {
	Recognized bool `json:"recognized"`
}

type LoginRequest struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

type LoginResponse struct {
	User  UserDTO `json:"user"`
	Token string  `json:"token"`
}

type MeResponse struct {
	User UserDTO `json:"user"`
}

type AddressDTO struct {
	Label   string `json:"label"`
	Line1   string `json:"line1"`
	Line2   string `json:"line2"`
	City    string `json:"city"`
	State   string `json:"state"`
	Pincode string `json:"pincode"`
}

func toAddressDTO(a entity.Address) AddressDTO {
	return AddressDTO{Label: a.Label, Line1: a.Line1, Line2: a.Line2, City: a.City, State: a.State, Pincode: a.Pincode}
}

func (a AddressDTO) toEntity() entity.Address {
	return entity.Address{Label: a.Label, Line1: a.Line1, Line2: a.Line2, City: a.City, State: a.State, Pincode: a.Pincode}
}

type CheckoutRequest struct {
	Email   string     `json:"email"`
	Phone   string     `json:"phone"`
	Address AddressDTO `json:"address"`
}

type CheckoutResponse struct {
	ID        int        `json:"id"`
	UserID    *int       `json:"user_id"` // null for guest checkouts
	Email     string     `json:"email"`
	Phone     string     `json:"phone"`
	Address   AddressDTO `json:"address"`
	CreatedAt time.Time  `json:"created_at"`
}

func toCheckoutResponse(c entity.Checkout) CheckoutResponse {
	return CheckoutResponse{
		ID:        c.ID,
		UserID:    c.UserID,
		Email:     c.Email,
		Phone:     c.Phone,
		Address:   toAddressDTO(c.Address),
		CreatedAt: c.CreatedAt,
	}
}

type SavedDetailsDTO struct {
	Phone      string     `json:"phone"`
	Address    AddressDTO `json:"address"`
	LastUsedAt time.Time  `json:"last_used_at"`
}

type SavedDetailsResponse struct {
	Saved []SavedDetailsDTO `json:"saved"` // most recently used first; empty for a first order
}
