package service

import (
	"context"
	"homework/internal/models"
	"homework/pkg/errors"
)

// helper function to convert pq.Int64Array to []int.
func GetWinCombSlice(pqArray []int) []int {
	winningNumbersInt := make([]int, len(pqArray))

	for i, val64 := range pqArray {
		winningNumbersInt[i] = int(val64)
	}

	return winningNumbersInt
}

func countMatches(ticketNumbers []int, winningNumbers []int) int {
	// Create a map of winning numbers for quick lookups
	winningSet := make(map[int]struct{}, len(winningNumbers))
	for _, n := range winningNumbers {
		winningSet[n] = struct{}{}
	}

	matchCount := 0
	// Use a map for ticket numbers to avoid double counting if ticket has duplicates (shouldn't happen ideally)
	ticketSet := make(map[int]struct{}, len(ticketNumbers))
	for _, num := range ticketNumbers {
		if _, alreadyChecked := ticketSet[num]; !alreadyChecked {
			if _, found := winningSet[num]; found {
				matchCount++
			}
			ticketSet[num] = struct{}{} // Mark as checked
		}
	}

	return matchCount
}

func ProcessTicket(ctx context.Context, ticket *models.TicketStore, repo Repository) (*models.TicketResult, error) {
	drawRes, err := repo.GetDraw(ctx, ticket.DrawId)
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

	ticketCombination, err := models.ParseTicketCombination(ticket.Data)
	if err != nil {
		return nil, errors.Errorf("couldn't parse ticket info")
	}

	result := countMatches(ticketCombination, drawWinCombination)

	// return fmt.Sprintf("combination here: %d, ticket combination: %w", result, ticketCombination), nil
	return &models.TicketResult{
		WinCombination: drawWinCombination,
		Combination:    ticketCombination,
		WinCount:       result,
		Id:             ticket.Id,
	}, nil
}
