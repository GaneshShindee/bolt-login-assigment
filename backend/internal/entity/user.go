// Package entity holds the domain models shared by every layer.
// They carry no JSON or database tags: each outer layer maps them to its own format.
package entity

type User struct {
	ID            int
	Email         string // always lowercase
	FirstName     string
	LastName      string
	LoginCodeHash string
}
