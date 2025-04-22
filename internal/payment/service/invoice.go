package service

import (
	"encoding/json"
	"homework/internal/models"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// Регистрарация инвойса
func RegisterInvoce(w http.ResponseWriter, r *http.Request) {
	var invoice *models.Invoice
	err := json.NewDecoder(r.Body).Decode(&invoice)
	if err != nil {
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
