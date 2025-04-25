package repository

import (
	"context"
	"homework/internal/models"
)

type Repository struct{}

func NewRepository() (*Repository, error) {
	return &Repository{}, nil
}

func (r *Repository) Create(ctx context.Context, ticket *models.Ticket) (*models.Ticket, error) {
	ticket.ID = "test-id-123"
	return ticket, nil
}

func (r *Repository) Get(ctx context.Context, id string) (*models.Ticket, error) {
	return &models.Ticket{
		ID:            id,
		UserID:        "1",
		DrawID:        "1",
		Numbers:       []int{1, 2, 3, 4, 5},
		Status:        models.TicketStatusActive,
		IsAutoNumbers: false,
	}, nil
}

func (r *Repository) GetByUserID(ctx context.Context, userID int64) ([]*models.Ticket, error) {
	return []*models.Ticket{
		{
			ID:            "1",
			UserID:        "1",
			DrawID:        "1",
			Numbers:       []int{1, 2, 3, 4, 5},
			Status:        models.TicketStatusActive,
			IsAutoNumbers: false,
		},
	}, nil
} 