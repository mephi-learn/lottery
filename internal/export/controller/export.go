package controller

import (
	"homework/internal/helpers"
	"net/http"
)

func (h *handler) ExportDraws(w http.ResponseWriter, r *http.Request) {
	draws, err := h.service.ExportDraws(r.Context())
	if err != nil {
		helpers.ErrorMessage(w, "export error", http.StatusBadRequest, err)
		return
	}

	helpers.SuccessMessage(w, "data", draws)
}
