package controller

import (
	"encoding/json"
	"fmt"
	"homework/internal/helpers"
	"homework/internal/models"
	"net/http"
	"strconv"
)

func (h *handler) RegisterPayment(w http.ResponseWriter, r *http.Request) {
	invoiceId, err := strconv.Atoi(r.PathValue("invoice_id"))
	if err != nil {
		helpers.ErrorMessage(w, fmt.Sprintf("invalid invoice: %s", r.PathValue("invoice_id")), http.StatusBadRequest, nil)
		return
	}

	paymentRequest := models.PaymentRequest{InvoiceID: invoiceId}
	err = json.NewDecoder(r.Body).Decode(&paymentRequest)
	if err != nil {
		helpers.ErrorMessage(w, "invalid json data", http.StatusBadRequest, err)
		return
	}

	// Регистрируем и проводим платёж
	if err = h.service.RegisterPayment(r.Context(), &paymentRequest); err != nil {
		helpers.ErrorMessage(w, "failed payment", http.StatusBadRequest, err)
		return
	}

	helpers.SuccessMessage(w, "invoice has been paid", nil)
}
