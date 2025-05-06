package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(invoiceId)
}
