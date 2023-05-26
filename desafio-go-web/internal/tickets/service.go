package tickets

import (
	"context"
)

type Service interface {
	GetTotalTickets(ctx context.Context, destination string) (int, error)
	AverageDestination(ctx context.Context, destination string) (float64, error)
}

type service struct {
	Storage Repository
}

func NewService(rp Repository) *service {
	return &service{
		Storage: rp,
	}
}
func (s *service) GetTotalTickets(ctx context.Context, destination string) (int, error) {
	tickets, err := s.Storage.GetTicketByDestination(ctx, destination)
	if err != nil {
		return 0, err
	}
	return len(tickets), nil
}
func (s *service) AverageDestination(ctx context.Context, destination string) (float64, error) {
	ticketsByDestination, err := s.Storage.GetTicketByDestination(ctx, destination)
	if err != nil {
		return 0, err
	}
	tickets, err := s.Storage.GetAll(ctx)
	if err != nil {
		return 0, err
	}
	average := float64(len(ticketsByDestination)) / float64(len(tickets))
	return average, nil
}
