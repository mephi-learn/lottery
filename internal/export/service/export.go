package service

import (
	"context"
	"homework/internal/models"
	"homework/pkg/errors"
)

func (s *exportService) ExportDraws(ctx context.Context) (*models.DrawExportResults, error) {
	draws, err := s.result.GetCompletedDraws(ctx)
	if err != nil {
		return nil, errors.Errorf("failed to get completed draw: %w", err)
	}

	results := &models.DrawExportResults{Draws: make([]*models.DrawExportResult, len(draws))}
	for i, draw := range draws {
		stat, err := s.draw.Drawing(ctx, draw.DrawId, draw.WinCombination)
		if err != nil {
			return nil, errors.Errorf("failed to get draw statistic: %w", err)
		}
		result := &models.DrawExportResult{
			DrawId:         draw.DrawId,
			WinCombination: draw.WinCombination,
			Statistic:      stat.Statistic,
			Tickets:        make(map[string][]int),
		}

		for resultName, drawTickets := range stat.WinTickets {
			tickets := make([]int, len(drawTickets))
			for j, ticket := range drawTickets {
				tickets[j] = ticket.Id
			}
			result.Tickets[resultName] = tickets
		}

		results.Draws[i] = result
	}

	return results, nil
}
