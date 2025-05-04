package controller

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"homework/internal/models"
	"net/http"
	"strconv"
	"time"
)

// RegisterInvoice регистрация инвойса.
func (h *handler) RegisterInvoice(w http.ResponseWriter, r *http.Request) {
	var invoice *models.Invoice
	if err := json.NewDecoder(r.Body).Decode(&invoice); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ticketId, err := strconv.Atoi(r.PathValue("ticketId"))
	if err != nil {
		http.Error(w, fmt.Sprintf("invalid ticket: %s", r.PathValue("ticketId")), http.StatusBadRequest)
		return
	}

	if err := h.service.RegisterInvoice(r.Context(), ticketId); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	invoice.ID = uuid.New()
	invoice.RegisterTime = time.Now()
	invoice.Status = "pending" // Начальный статус

	// invoices[invoice.ID] = invoice - данные с бэка

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(invoice)
}
