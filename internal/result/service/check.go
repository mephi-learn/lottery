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

	drawRes, err := s.repo.GetDraw(ctx, ticket.DrawId)

	if err != nil {
		return nil, errors.Errorf("failed to get draw: %w", err)
	}
	if drawRes == nil {
		return nil, errors.Errorf("draw not found")
	}

	// compare ticket numbers with winning combination
	if drawRes.WinCombination == nil {
		return nil, errors.Errorf("winning combination not found")
	}

	drawWinCombination := GetWinCombSlice(drawRes.WinCombination)

	ticketCombination, err := ParseTicketCombination(ticket.Data)

	if err != nil {
		return nil, errors.Errorf("couldn't parse ticket info")
	}

	result := countMatches(ticketCombination, drawWinCombination)

	// return fmt.Sprintf("combination here: %d, ticket combination: %w", result, ticketCombination), nil
	return &models.TicketResult{
		WinCombination: drawWinCombination,
		Combination:    ticketCombination,
		WinCount:       result,
	}, nil
}
