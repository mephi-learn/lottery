package repository

import (
	"context"
	"homework/internal/models"
	"homework/pkg/errors"
)

// DebitingFundsFromWallet Списание средств с кошелька пользователя
func (r *repository) DebitingFundsFromWallet(ctx context.Context, amount float64) error {
	user, err := models.UserFromContext(ctx)
	if err != nil {
		return errors.Errorf("authentificate need: %w", err)
	}

	_, err = r.db.ExecContext(ctx, "update users set wallet = wallet - $1 where id = $2", amount, user.ID)
	if err != nil {
		return errors.Errorf("failed to debiting funds from wallet: %w", err)
	}

	return nil
}

// GetAmountInUserWallet Получение суммы на кошельке пользователя
func (r *repository) GetAmountInUserWallet(ctx context.Context) (float64, error) {
	user, err := models.UserFromContext(ctx)
	if err != nil {
		return -1, errors.Errorf("authentificate need: %w", err)
	}

	var amountInWallet float64

	if err := r.db.QueryRowContext(ctx, "select wallet from users where id = $1", user.ID).Scan(&amountInWallet); err != nil {
		return -1, errors.Errorf("failed to debiting funds from wallet: %w", err)
	}

	return amountInWallet, nil
}
