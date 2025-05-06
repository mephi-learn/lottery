package service

import (
	"context"
)

type Service struct{}

func New() (*Service, error) {
	return &Service{}, nil
}

type DrawExport struct {
	DrawID             int64
	LotteryType        string
	WinningCombination []int
	WinnerCount        int
}

func (s *Service) ExportDraws(ctx context.Context) ([]DrawExport, error) {
	return []DrawExport{
		{
			DrawID:             1,
			LotteryType:        "5 из 36",
			WinningCombination: []int{1, 2, 3, 4, 5},
			WinnerCount:        10,
		},
	}, nil
}
