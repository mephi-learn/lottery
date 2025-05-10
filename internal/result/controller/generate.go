package controller

import (
	"fmt"
	"homework/internal/helpers"
	"net/http"
	"strconv"
)

func (h *handler) Drawing(w http.ResponseWriter, r *http.Request) {
	// Парсим входные данные
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		helpers.ErrorMessage(w, fmt.Sprintf("invalid id: %s", r.PathValue("id")), http.StatusBadRequest, err)
		return
	}

	result, err := h.service.Drawing(r.Context(), id)
	if err != nil {
		h.log.Error("failed to generate draw results", "err", err)
		helpers.ErrorMessage(w, "failed to get draw results", http.StatusBadRequest, err)

		return
	}

	helpers.SuccessMessage(w, "result", result)
}
