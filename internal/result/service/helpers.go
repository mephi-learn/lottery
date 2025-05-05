package service

import (
	"context"
	"encoding/base64"
	"fmt"
	"homework/internal/models"
	"homework/pkg/errors"
	"strconv"
	"strings"

	"github.com/lib/pq"
)

func ParseTicketCombination(combination string) (ticketNumbers []int, err error) {
	data, err := base64.StdEncoding.DecodeString(combination)
	if err != nil {
		return nil, errors.New("unknown decode ticket data")
	}
	parts := strings.SplitN(string(data), ";", 2)

	numberStr := parts[1]

	digitStrings := strings.Split(numberStr, ",")
	ticketNumbers = make([]int, 0, len(digitStrings))

	for i, digitStr := range digitStrings {
		digitStr = strings.TrimSpace(digitStr) // Handle potential spaces
		if digitStr == "" {
			err = fmt.Errorf("invalid ticket combination format: empty number string at index %d", i)
			return
		}

		digit, parseErr := strconv.Atoi(digitStr)
		if parseErr != nil {
			err = fmt.Errorf("invalid ticket combination format: failed to parse number '%s' at index %d: %w", digitStr, i, parseErr)
			return
		}
		ticketNumbers = append(ticketNumbers, digit)
	}

	return ticketNumbers, nil
}

// helper function to convert pq.Int64Array to []int
func GetWinCombSlice(pqArray pq.Int64Array) []int {
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
		Id:            ticket.Id,
	}, nil
}