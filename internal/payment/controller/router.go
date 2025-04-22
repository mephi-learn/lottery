package controller

import (
	"context"
	"homework/pkg/log"
	"net/http"
	"time"
)

type handler struct {
	service paymentService
	log     log.Logger
}

type paymentService interface {
	RegisterInvoice(ctx context.Context, timeRegistration time.Time) (err error)
	RegisterPayment(ctx context.Context, payment float64) (err error)
}

func (h *handler) WithRouter(mux *http.ServeMux) {
	// Invoice
	mux.Handle("POST /api/invoice", http.HandlerFunc(h.RegisterInvoice))
	// Payment
	mux.Handle("POST /api/payments", http.HandlerFunc(h.RegisterPayment))
}
