package models

// User represents a user purchasing a ticket
type User struct {
	FirstName string
	LastName  string
	Email     string
}

// Ticket represents a purchased train ticket
type Ticket struct {
	ReceiptID string
	From      string
	To        string
	User      User
	Seat      string
	PricePaid float64
}
