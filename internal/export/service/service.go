package service

import (
	"context"
	"homework/internal/models"
	"homework/pkg/errors"
	"homework/pkg/log"
)

// Repository реализует интерфейс репозитория сервиса покупки билетов.
type Repository interface{}

type DrawService interface {
	ListCompletedDraw(ctx context.Context) ([]*models.DrawOutput, error)
	Drawing(ctx context.Context, drawId int, combination []int) (*models.DrawingResult, error)
}

type ResultService interface {
	GetCompletedDraws(ctx context.Context) ([]*models.DrawResultStore, error)
}

type ExportOption func(*exportService) error

type exportService struct {
	repo Repository

	draw   DrawService
	result ResultService

	log log.Logger
}

// NewExportService возвращает имплементацию сервиса для вывода статистики.
func NewExportService(opts ...ExportOption) (*exportService, error) {
	var svc exportService

	for _, opt := range opts {
		opt(&svc)
	}

	if svc.log == nil {
		return nil, errors.Errorf("no logger provided")
	}

	if svc.draw == nil {
		return nil, errors.Errorf("no draw provided")
	}

	return &svc, nil
}

func WithExportLogger(logger log.Logger) ExportOption {
	return func(r *exportService) error {
		r.log = logger
		return nil
	}
}

func WithExportRepository(repo Repository) ExportOption {
	return func(r *exportService) error {
		r.repo = repo
		return nil
	}
}

func WithDrawService(draw DrawService) ExportOption {
	return func(r *exportService) error {
		r.draw = draw
		return nil
	}
}

func WithResultService(result ResultService) ExportOption {
	return func(r *exportService) error {
		r.result = result
		return nil
	}
}
