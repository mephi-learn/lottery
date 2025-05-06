package service

import (
	"context"
	"homework/internal/models"
	"homework/pkg/errors"
)

func (s *resultService) CheckTicketResult(ctx context.Context, ticketId, userId int) (*models.TicketResult, error) {
	ticket, err := s.repo.GetUserTicket(ctx, ticketId, userId)
	if err != nil {
		return nil, errors.Errorf("failed to get ticket: %w", err)
	}

	if ticket == nil {
		return nil, errors.Errorf("ticket not found")
	}

	return ProcessTicket(ctx, ticket, s.repo)
}

func (s *resultService) CheckTicketsResult(ctx context.Context, userId int) ([]models.TicketResult, error) {
	tickets, err := s.repo.GetUserTickets(ctx, userId)
	if err != nil {
		return nil, errors.Errorf("failed to get ticket: %w", err)
	}

	if len(tickets) == 0 {
		return nil, errors.Errorf("tickets not found")
	}

	resTickets := []models.TicketResult{}

	for _, ticket := range tickets {
		res, err := ProcessTicket(ctx, &ticket, s.repo)
		if err != nil {
			return nil, errors.Errorf("ticket processing error")
		}
		if res != nil {
			resTickets = append(resTickets, *res)
		}
	}

	return resTickets, nil
}
