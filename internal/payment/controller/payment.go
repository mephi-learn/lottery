package controller

import (
	"encoding/json"
	"fmt"
	"homework/internal/models"
	"net/http"
	"strconv"
)

func (h *handler) RegisterPayment(w http.ResponseWriter, r *http.Request) {
	invoiceId, err := strconv.Atoi(r.PathValue("invoice_id"))
	if err != nil {
		http.Error(w, fmt.Sprintf("invalid invoice: %s", r.PathValue("invoice_id")), http.StatusBadRequest)
		return
	}

	paymentRequest := models.PaymentRequest{InvoiceID: invoiceId}
	err = json.NewDecoder(r.Body).Decode(&paymentRequest)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed decore request: %s", err.Error()), http.StatusBadRequest)
		return
	}

	// Регистрируем и проводим платёж
	if err = h.service.RegisterPayment(r.Context(), &paymentRequest); err != nil {
		http.Error(w, fmt.Sprintf("failed payment: %s", err.Error()), http.StatusBadRequest)
		return
	}

	// mockPaymentData := 1000.5 // просто временный мок
	//h.service.RegisterPayment(request.Context(), mockPaymentData)
	// invoice_id
}
