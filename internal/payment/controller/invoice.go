package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type response struct {
	Message   string `json:"message"`
	InvoiceId int    `json:"invoice_id"`
}

// RegisterInvoice регистрация инвойса. (по ticketid - создается инвойс текущему пользователю)
func (h *handler) RegisterInvoice(w http.ResponseWriter, r *http.Request) {
	ticketId, err := strconv.Atoi(r.PathValue("ticket_id"))
	if err != nil {
		http.Error(w, fmt.Sprintf("invalid ticket: %s", r.PathValue("ticket_id")), http.StatusBadRequest)
		return
	}

	invoiceId, err := h.service.RegisterInvoice(r.Context(), ticketId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := response{
		Message:   "invoice has been created",
		InvoiceId: invoiceId,
	}

	result, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed create response: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(result)
}
