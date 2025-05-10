package controller

import (
	"fmt"
	"homework/internal/helpers"
	"net/http"
	"strconv"
)

func (h *handler) GetDrawResults(w http.ResponseWriter, r *http.Request) {
	// Парсим входные данные
	id, err := strconv.Atoi(r.PathValue("draw_id"))
	if err != nil {
		helpers.ErrorMessage(w, fmt.Sprintf("invalid draw id: %s", r.PathValue("draw_id")), http.StatusBadRequest, err)
		return
	}

	result, err := h.service.GetDrawResults(r.Context(), id)
	if err != nil {
		h.log.Error("failed to get draw results", "err", err)
		helpers.ErrorMessage(w, "failed to get draw results", http.StatusBadRequest, err)

		return
	}

	helpers.SuccessMessage(w, "result", result)
}
