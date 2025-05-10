package controller

import (
	"fmt"
	"homework/internal/helpers"
	"net/http"
	"strconv"
)

// RegisterInvoice регистрация инвойса. (по ticketid - создается инвойс текущему пользователю).
func (h *handler) RegisterInvoice(w http.ResponseWriter, r *http.Request) {
	ticketId, err := strconv.Atoi(r.PathValue("ticket_id"))
	if err != nil {
		helpers.ErrorMessage(w, fmt.Sprintf("invalid ticket: %s", r.PathValue("ticket_id")), http.StatusBadRequest, nil)
		return
	}

	invoiceId, err := h.service.RegisterInvoice(r.Context(), ticketId)
	if err != nil {
		helpers.ErrorMessage(w, "failed register invoice", http.StatusBadRequest, err)
		return
	}

	helpers.SuccessMessage(w, "invoice create", map[string]any{"invoice id": invoiceId})
}
