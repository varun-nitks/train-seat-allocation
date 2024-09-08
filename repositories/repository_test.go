package repositories

import (
	"testing"
	"train-seat-allocation/models"
)

func TestAddTicket(t *testing.T) {
	repo := NewInMemoryTicketRepository()
	ticket := &models.Ticket{
		User: models.User{Email: "user1@example.com"},
		Seat: "A1",
	}
	err := repo.AddTicket(ticket)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestGetTicket(t *testing.T) {
	repo := NewInMemoryTicketRepository()
	ticket := &models.Ticket{
		User: models.User{Email: "user1@example.com"},
		Seat: "A1",
	}
	repo.AddTicket(ticket)

	tests := []struct {
		email    string
		expected *models.Ticket
		errMsg   string
	}{
		{"user1@example.com", ticket, ""},
		{"user2@example.com", nil, "ticket not found"},
	}

	for _, tt := range tests {
		ticket, err := repo.GetTicket(tt.email)
		if err != nil && err.Error() != tt.errMsg {
			t.Errorf("expected error: %v, got: %v", tt.errMsg, err)
		}
		if ticket != nil && ticket.Seat != tt.expected.Seat {
			t.Errorf("expected seat: %v, got: %v", tt.expected.Seat, ticket.Seat)
		}
	}
}

func TestRemoveTicket(t *testing.T) {
	repo := NewInMemoryTicketRepository()
	ticket := &models.Ticket{
		User: models.User{Email: "user1@example.com"},
		Seat: "A1",
	}
	repo.AddTicket(ticket)

	err := repo.RemoveTicket("user1@example.com")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	err = repo.RemoveTicket("user2@example.com")
	if err == nil || err.Error() != "ticket not found" {
		t.Errorf("expected error: ticket not found, got %v", err)
	}
}

func TestIsSeatAllocated(t *testing.T) {
	repo := NewInMemoryTicketRepository()
	repo.InitializeSeatAllocations("A1", "A")

	tests := []struct {
		seat      string
		section   string
		allocated bool
	}{
		{"A1", "A", true},
		{"A2", "A", false},
		{"A1", "B", false},
	}

	for _, tt := range tests {
		allocated := repo.IsSeatAllocated(tt.seat, tt.section)
		if allocated != tt.allocated {
			t.Errorf("expected %v, got %v", tt.allocated, allocated)
		}
	}
}

func TestPurchaseTicket(t *testing.T) {
	repo := NewInMemoryTicketRepository()
	user := models.User{Email: "user1@example.com"}
	ticket, err := repo.PurchaseTicket(user, "StationA", "StationB")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if ticket.User.Email != user.Email {
		t.Errorf("expected email: %v, got: %v", user.Email, ticket.User.Email)
	}
}

func TestModifySeat(t *testing.T) {
	repo := NewInMemoryTicketRepository()
	repo.AddTicket(&models.Ticket{
		User: models.User{Email: "user1@example.com"},
		Seat: "A1",
	})
	repo.AddTicket(&models.Ticket{
		User: models.User{Email: "user2@example.com"},
		Seat: "A2",
	})

	tests := []struct {
		email    string
		newSeat  string
		expected string
		errMsg   string
	}{
		{"user1@example.com", "A2", "A1", "seat already allocated to another passenger"},
		{"user3@example.com", "A3", "A3", ""},
		{"user1@example.com", "A4", "A4", ""},
	}

	for _, tt := range tests {
		err := repo.ModifySeat(tt.email, tt.newSeat)
		if err != nil && err.Error() != tt.errMsg {
			t.Errorf("expected error: %v, got: %v", tt.errMsg, err)
		}
		if err == nil {
			ticket, _ := repo.GetTicket(tt.email)
			if ticket.Seat != tt.expected {
				t.Errorf("expected seat: %v, got: %v", tt.expected, ticket.Seat)
			}
		}
	}
}

func TestGetUsersBySection(t *testing.T) {
	repo := NewInMemoryTicketRepository()
	repo.AddTicket(&models.Ticket{
		User: models.User{Email: "user1@example.com"},
		Seat: "A1",
	})
	repo.AddTicket(&models.Ticket{
		User: models.User{Email: "user2@example.com"},
		Seat: "B1",
	})
	repo.InitializeSeatAllocations("A1", "A")
	repo.InitializeSeatAllocations("B1", "B")

	tests := []struct {
		section  string
		expected int
	}{
		{"A", 1},
		{"B", 1},
		{"C", 0},
	}

	for _, tt := range tests {
		users, err := repo.GetUsersBySection(tt.section)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if len(users) != tt.expected {
			t.Errorf("expected %v users, got %v", tt.expected, len(users))
		}
	}
}
