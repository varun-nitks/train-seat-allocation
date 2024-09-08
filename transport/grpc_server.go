package transport

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"train-seat-allocation/models"
	pb "train-seat-allocation/proto/pb"
	"train-seat-allocation/services"
)

type GRPCServer struct {
	ticketService *services.TicketService
	pb.UnimplementedTrainTicketServiceServer
}

func NewGRPCServer(service *services.TicketService) *GRPCServer {
	return &GRPCServer{ticketService: service}
}

func (s *GRPCServer) PurchaseTicket(ctx context.Context, req *pb.PurchaseTicketRequest) (*pb.PurchaseTicketResponse, error) {
	user := models.User{
		FirstName: req.User.FirstName,
		LastName:  req.User.LastName,
		Email:     req.User.Email,
	}

	ticket, err := s.ticketService.PurchaseTicket(user, req.From, req.To)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to purchase ticket: %v", err)
	}

	return &pb.PurchaseTicketResponse{
		ReceiptId: ticket.ReceiptID,
		From:      ticket.From,
		To:        ticket.To,
		User:      req.User,
		Seat:      ticket.Seat,
		PricePaid: ticket.PricePaid,
	}, nil
}

func (s *GRPCServer) ModifyUserSeat(ctx context.Context, req *pb.ModifyUserSeatRequest) (*pb.ModifyUserSeatResponse, error) {
	err := s.ticketService.ModifySeat(req.Email, req.NewSeat)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to modify seat: %v", err)
	}

	return &pb.ModifyUserSeatResponse{
		Success: true,
	}, nil
}

func (s *GRPCServer) GetReceipt(ctx context.Context, req *pb.GetReceiptRequest) (*pb.GetReceiptResponse, error) {
	ticket, err := s.ticketService.GetTicket(req.Email)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "failed to get receipt: %v", err)
	}

	return &pb.GetReceiptResponse{
		ReceiptId: ticket.ReceiptID,
		From:      ticket.From,
		To:        ticket.To,
		User: &pb.User{
			FirstName: ticket.User.FirstName,
			LastName:  ticket.User.LastName,
			Email:     ticket.User.Email,
		},
		Seat:      ticket.Seat,
		PricePaid: ticket.PricePaid,
	}, nil
}

func (s *GRPCServer) RemoveUser(ctx context.Context, req *pb.RemoveUserRequest) (*pb.RemoveUserResponse, error) {
	err := s.ticketService.RemoveTicket(req.Email)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "failed to remove user: %v", err)
	}

	return &pb.RemoveUserResponse{
		Success: true,
	}, nil
}

func (s *GRPCServer) GetUsersBySection(ctx context.Context, req *pb.GetUsersBySectionRequest) (*pb.GetUsersBySectionResponse, error) {
	seatAllocations, err := s.ticketService.GetUsersBySection(req.Section)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get users by section: %v", err)
	}

	return &pb.GetUsersBySectionResponse{
		Users: seatAllocations,
	}, nil
}
