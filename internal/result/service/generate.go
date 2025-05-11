package service

import (
	"context"
	"homework/internal/models"
	"homework/pkg/errors"
	"time"
)

func (s *resultService) Drawing(ctx context.Context, drawId int) ([]int, error) {
	draw, err := s.repo.GetDraw(ctx, drawId)
	if err != nil {
		return nil, errors.Errorf("failed to get draw: %w", err)
	}

	if draw == nil {
		return nil, errors.Errorf("draw not found")
	}

	// Если тираж не в состоянии "запланирован", выдаём ошибку
	if draw.DrawStatusId != int(models.DrawStatusPlanned) {
		return nil, errors.Errorf("draw not completed")
	}

	// Переводим тираж в состояние active (розыгрыш в процессе)
	if err = s.draw.ActiveDraw(ctx, drawId); err != nil {
		return nil, errors.Errorf("failed change status to active: %w", err)
	}

	lottery, err := s.lottery.LotteryByType(draw.LotteryType)
	if err != nil {
		_ = s.draw.FailedDraw(ctx, drawId)
		return nil, errors.Errorf("failed to get lottery: %w", err)
	}

	// Generate the winning numbers based on the lottery type
	winningNumbers, err := lottery.GenerateWinningCombination()
	if err != nil {
		_ = s.draw.FailedDraw(ctx, drawId)
		return nil, errors.Errorf("failed to generate winning numbers: %w", err)
	}

	// save the draw result with the winning numbers
	err = s.repo.SaveWinCombination(ctx, drawId, winningNumbers)
	if err != nil {
		_ = s.draw.FailedDraw(ctx, drawId)
		return nil, errors.Errorf("failed to save winning numbers: %w", err)
	}

	_ = s.draw.CompletedDraw(ctx, drawId)

	// Маркируем билеты
	_, err = s.draw.DrawingAndMarkTickets(ctx, drawId, winningNumbers)
	if err != nil {

	}

	return winningNumbers, nil
}

func (s *resultService) StartDrawsReadyToBegin(ctx context.Context) {
	s.log.InfoContext(ctx, "starting search draw ready to begin")
	ticker := time.NewTicker(time.Minute)
	go func() {
		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				s.log.InfoContext(ctx, "search draw ready to begin stopped")

				return
			case <-ticker.C:
				s.log.InfoContext(ctx, "search draw ready to begin")
				draws, err := s.draw.GetReadyToBeginDraws(ctx)
				if err != nil {
					s.log.ErrorContext(ctx, "failed to draws", "error", err)
					continue
				}

				if len(draws) > 0 {
					s.log.InfoContext(ctx, "found ready to begin draws", "count", len(draws))
				}

				for _, draw := range draws {
					if _, err = s.Drawing(ctx, draw.Id); err != nil {
						s.log.ErrorContext(ctx, "failed to drawing", "draw_id", draw.Id, "error", err)
					}
				}
			}
		}
	}()
}
