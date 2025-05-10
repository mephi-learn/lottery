package controller

import (
	"encoding/json"
	"fmt"
	"homework/internal/helpers"
	"net/http"
	"strconv"
)

type Combination struct {
	Combination []int `json:"combination"`
}

// RegisterCustomInvoice регистрация инвойса с билетом.
func (h *handler) RegisterCustomInvoice(w http.ResponseWriter, r *http.Request) {
	drawId, err := strconv.Atoi(r.PathValue("draw_id"))
	if err != nil {
		helpers.ErrorMessage(w, fmt.Sprintf("invalid draw: %s", r.PathValue("draw_id")), http.StatusBadRequest, nil)
		return
	}

	var combination Combination
	if err := json.NewDecoder(r.Body).Decode(&combination); err != nil {
		helpers.ErrorMessage(w, "invalid json data", http.StatusBadRequest, err)
		return
	}

	invoiceId, err := h.service.RegisterCustomInvoice(r.Context(), drawId, combination.Combination)
	if err != nil {
		helpers.ErrorMessage(w, "failed register invoice", http.StatusBadRequest, err)
		return
	}

	helpers.SuccessMessage(w, "invoice create", map[string]any{"invoice id": invoiceId})
}
