package repositories

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"train-seat-allocation/models"
)

// TicketRepository interface defines methods for ticket operations
type TicketRepository interface {
	InitializeSeatAllocations(seat string, section string)
	AddTicket(ticket *models.Ticket) error
	GetTicket(email string) (*models.Ticket, error)
	RemoveTicket(email string) error
	IsSeatAllocated(seat, section string) bool
	AllocateSeat() string
	PurchaseTicket(user models.User, from, to string) (*models.Ticket, error)
	ModifySeat(email, newSeat string) error
	GetUsersBySection(section string) ([]*models.Ticket, error)
}

// InMemoryTicketRepository is an in-memory implementation of TicketRepository
type InMemoryTicketRepository struct {
	tickets         map[string]*models.Ticket
	seatAllocations map[string]string
	mu              sync.Mutex
}

// NewInMemoryTicketRepository creates a new in-memory ticket repository
func NewInMemoryTicketRepository() TicketRepository {
	return &InMemoryTicketRepository{
		tickets:         make(map[string]*models.Ticket),
		seatAllocations: make(map[string]string),
	}
}

// InitializeSeatAllocations Add a method to initialize seat allocations in the InMemoryTicketRepository struct
func (r *InMemoryTicketRepository) InitializeSeatAllocations(seat string, section string) {
	r.seatAllocations[seat] = section
}

func (r *InMemoryTicketRepository) AddTicket(ticket *models.Ticket) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tickets[ticket.User.Email] = ticket
	return nil
}

func (r *InMemoryTicketRepository) GetTicket(email string) (*models.Ticket, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	ticket, exists := r.tickets[email]
	if !exists {
		return nil, errors.New("ticket not found")
	}
	return ticket, nil
}

func (r *InMemoryTicketRepository) RemoveTicket(email string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	ticket, exists := r.tickets[email]
	if !exists {
		return errors.New("ticket not found")
	}

	// Deallocate the seat
	delete(r.seatAllocations, ticket.Seat)

	// Remove the ticket
	delete(r.tickets, email)
	return nil
}

// IsSeatAllocated checks if a seat is already allocated
func (r *InMemoryTicketRepository) IsSeatAllocated(seat, section string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	allocatedSection, exists := r.seatAllocations[seat]
	return exists && allocatedSection == section
}

// AllocateSeat marks a seat as allocated
func (r *InMemoryTicketRepository) AllocateSeat() string {
	var seat, section string
	sections := []string{"A", "B"}
	for {
		section = sections[rand.Intn(len(sections))]
		seat = fmt.Sprintf("%s%d", section, rand.Intn(100))
		if !r.IsSeatAllocated(seat, section) {
			r.mu.Lock()
			r.InitializeSeatAllocations(seat, section)
			r.mu.Unlock()
			return seat
		}
	}
}

// PurchaseTicket handles ticket purchase
func (r *InMemoryTicketRepository) PurchaseTicket(user models.User, from, to string) (*models.Ticket, error) {
	seat := r.AllocateSeat()
	ticket := &models.Ticket{
		ReceiptID: fmt.Sprintf("R%d", rand.Intn(10000)),
		From:      from,
		To:        to,
		User:      user,
		Seat:      seat,
		PricePaid: 20.00,
	}
	err := r.AddTicket(ticket)
	return ticket, err
}

func (r *InMemoryTicketRepository) ModifySeat(email, newSeat string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if the new seat is already allocated
	for _, ticket := range r.tickets {
		if ticket.Seat == newSeat {
			return errors.New("seat already allocated to another passenger")
		}
	}

	ticket, exists := r.tickets[email]
	if !exists {
		// Allocate new seat if ticket not found
		r.tickets[email] = &models.Ticket{
			User: models.User{Email: email},
			Seat: newSeat,
		}
		return nil
	}

	ticket.Seat = newSeat
	return nil
}

func (r *InMemoryTicketRepository) GetUsersBySection(section string) ([]*models.Ticket, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var users []*models.Ticket
	for seat, sec := range r.seatAllocations {
		if sec == section {
			for _, ticket := range r.tickets {
				if ticket.Seat == seat {
					users = append(users, ticket)
				}
			}
		}
	}
	return users, nil
}
