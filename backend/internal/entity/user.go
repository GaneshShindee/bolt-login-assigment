package entity

type User struct {
	ID            int
	Email         string // always lowercase
	FirstName     string
	LastName      string
	LoginCodeHash string
}
