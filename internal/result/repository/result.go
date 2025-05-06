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

// GetCompletedDraw get draw status and it's winning combination (if there is one already)
func (r *repository) GetCompletedDraws(ctx context.Context) ([]*models.DrawResultStore, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT
			d.id,
			d.status_id,
			d.lottery_type,
			r.win_combination
		FROM draws d LEFT JOIN draw_results r ON r.draw_id = d.id
		WHERE d.status_id = $1`, models.DrawStatusCompleted)
	if err != nil {
		return nil, errors.Errorf("failed to get draw info: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	var draws []*models.DrawResultStore
	for rows.Next() {
		var draw models.DrawResultStore
		if err = rows.Scan(&draw.Id, &draw.DrawStatusId, &draw.LotteryType, &draw.WinCombination); err != nil {
			return draws, err
		}
		draws = append(draws, &draw)
	}
	if err = rows.Err(); err != nil {
		return draws, err
	}

	return draws, nil
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
		SELECT id, status_id, draw_id, data, user_id FROM tickets WHERE id = $1 AND user_id = $2`, ticketId, userId).Scan(
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

func (r *repository) GetUserTickets(ctx context.Context, userId int) ([]models.TicketStore, error) {
	query := `
		SELECT id, status_id, draw_id, data, user_id 
		FROM tickets 
		WHERE user_id = $1
		ORDER BY id DESC`

	rows, err := r.db.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, errors.Errorf("failed to query tickets for user_id %d, err: %w", userId, err)
	}
	defer rows.Close()

	tickets := []models.TicketStore{} // Initialize an empty slice

	for rows.Next() {
		var ticket models.TicketStore
		if err := rows.Scan(
			&ticket.Id,
			&ticket.StatusId,
			&ticket.DrawId,
			&ticket.Data, // Ensure models.TicketStore.Data field can handle the DB type
			&ticket.UserId,
		); err != nil {
			// Error scanning a specific row
			return nil, errors.Errorf("failed to scan ticket row, cause: %w", err)
		}
		tickets = append(tickets, ticket)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Errorf("error iterating over ticket rows, cause: %w", err)
	}

	return tickets, nil
}
