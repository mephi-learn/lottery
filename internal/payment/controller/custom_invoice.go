package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type Combination struct {
	Combination []int `json:"combination"`
}

// RegisterCustomInvoice регистрация инвойса с билетом
func (h *handler) RegisterCustomInvoice(w http.ResponseWriter, r *http.Request) {
	drawId, err := strconv.Atoi(r.PathValue("draw_id"))
	if err != nil {
		http.Error(w, fmt.Sprintf("invalid draw: %s", r.PathValue("draw_id")), http.StatusBadRequest)
		return
	}

	var combination Combination
	if err := json.NewDecoder(r.Body).Decode(&combination); err != nil {
		http.Error(w, fmt.Sprintf("failed decore request: %s", err.Error()), http.StatusBadRequest)
		return
	}

	invoiceId, err := h.service.RegisterCustomInvoice(r.Context(), drawId, combination.Combination)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := ResponseInvoice{
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
