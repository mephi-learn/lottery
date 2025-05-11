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

	// Регистрируем и проводим платёж
	if err = h.service.RegisterPayment(r.Context(), invoiceId); err != nil {
		helpers.ErrorMessage(w, "failed payment", http.StatusBadRequest, err)
		return
	}

	helpers.SuccessMessage(w, "invoice has been paid", nil)
}

func (h *handler) FillWallet(w http.ResponseWriter, r *http.Request) {
	var paymentRequest *models.PaymentRequest
	err := json.NewDecoder(r.Body).Decode(&paymentRequest)
	if err != nil {
		helpers.ErrorMessage(w, "invalid json data", http.StatusBadRequest, err)
		return
	}
	if paymentRequest.Price <= 0 {
		helpers.ErrorMessage(w, "invalid price", http.StatusBadRequest, err)
		return
	}

	// Пополняем кошелек пользователя
	if err = h.service.FillWallet(r.Context(), paymentRequest); err != nil {
		helpers.ErrorMessage(w, "failed funds transfer", http.StatusBadRequest, err)
		return
	}

	helpers.SuccessMessage(w, "the wallet has been replenished", nil)
}
