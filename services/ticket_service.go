package services

import (
	"train-seat-allocation/models"
	pb "train-seat-allocation/proto/pb"
	"train-seat-allocation/repositories"
)

type TicketService struct {
	ticketRepo repositories.TicketRepository
}

func NewTicketService(repo repositories.TicketRepository) *TicketService {
	return &TicketService{ticketRepo: repo}
}

func (s *TicketService) PurchaseTicket(user models.User, from, to string) (*models.Ticket, error) {
	return s.ticketRepo.PurchaseTicket(user, from, to)
}

func (s *TicketService) ModifySeat(email, newSeat string) error {
	return s.ticketRepo.ModifySeat(email, newSeat)
}

func (s *TicketService) GetTicket(email string) (*models.Ticket, error) {
	return s.ticketRepo.GetTicket(email)
}

func (s *TicketService) RemoveTicket(email string) error {
	return s.ticketRepo.RemoveTicket(email)
}

func (s *TicketService) GetUsersBySection(section string) ([]*pb.SeatAllocation, error) {
	users, err := s.ticketRepo.GetUsersBySection(section)
	if err != nil {
		return nil, err
	}

	var seatAllocations []*pb.SeatAllocation
	for _, user := range users {
		seatAllocations = append(seatAllocations, &pb.SeatAllocation{
			User: &pb.User{
				FirstName: user.User.FirstName,
				LastName:  user.User.LastName,
				Email:     user.User.Email,
			},
			Seat: user.Seat,
		})
	}
	return seatAllocations, nil
}
