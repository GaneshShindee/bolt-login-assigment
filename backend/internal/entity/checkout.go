package entity

import "time"

// Address is an Indian shipping address.
type Address struct {
	Label   string // "Home", "Work" or a custom name
	Line1   string // house / flat, street
	Line2   string // optional: area, landmark
	City    string
	State   string
	Pincode string
}

type Checkout struct {
	ID        int
	UserID    *int // nil for guest checkouts
	Email     string
	Phone     string
	Address   Address
	CreatedAt time.Time
}

// SavedDetails is a phone + address from a user's earlier order, offered again at checkout.
type SavedDetails struct {
	Phone      string
	Address    Address
	LastUsedAt time.Time
}
