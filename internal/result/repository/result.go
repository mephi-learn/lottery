package repository

import (
	"context"
	"database/sql"
	"homework/internal/models"
	"homework/pkg/errors"

	"github.com/lib/pq"
)

// get draw status and it's winning combination (if there is one already)
func (r *repository) GetDraw(ctx context.Context, drawId int) (*models.DrawResultStore, error) {
	drawRes := models.DrawResultStore{}
	if err := r.db.QueryRowContext(ctx, `
		SELECT
			d.id,
			d.status_id,
			d.lottery_type,
			r.win_combination
		FROM draws d
			LEFT JOIN draw_results r ON r.draw_id = d.id
		WHERE d.id = $1`, drawId).Scan(
		&drawRes.Id,
		&drawRes.DrawStatusId,
		&drawRes.LotteryType,
		&drawRes.WinCombination,
	); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Errorf("failed to get draw info: %w", err)
		}

		return nil, nil
	}

	return &drawRes, nil
}

func (r *repository) SaveWinCombination(ctx context.Context, drawId int, winCombination []int) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO draw_results (draw_id, win_combination)
		VALUES ($1, $2)
	`, drawId, pq.Array(winCombination))

	if err != nil {
		return errors.Errorf("failed to save winning combination: %w", err)
	}

	return nil
}

func (r *repository) GetUserTicket(ctx context.Context, ticketId, userId int) (*models.TicketStore, error) {
	ticket := models.TicketStore{}
	if err := r.db.QueryRowContext(ctx, `
		SELECT * FROM tickets WHERE id = $1 AND user_id = $2`, ticketId, userId).Scan(
		&ticket.Id,
		&ticket.StatusId,
		&ticket.DrawId,
		&ticket.Data,
		&ticket.UserId,
	); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Errorf("failed to get ticket info: %w", err)
		}

		return nil, nil
	}

	return &ticket, nil
}
