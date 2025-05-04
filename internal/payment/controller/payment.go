package controller

import (
	"fmt"
	"net/http"
	"strconv"
)

func (h *handler) RegisterPayment(w http.ResponseWriter, r *http.Request) {
	ticketId, err := strconv.Atoi(r.PathValue("invoice_id"))
	if err != nil {
		http.Error(w, fmt.Sprintf("invalid invoice: %s", r.PathValue("invoice_id")), http.StatusBadRequest)
		return
	}

	mockPaymentData := 1000.5 // просто временный мок
	h.service.RegisterPayment(request.Context(), mockPaymentData)
	invoice_id
}
