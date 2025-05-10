package controller

import (
	"homework/internal/helpers"
	"net/http"
)

func (h *handler) ListActiveDraws(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.log.InfoContext(ctx, "list active draws", "id", r.PathValue("draw_id"))

	// Получаем список активных тиражей
	list, err := h.service.ListActiveDraws(ctx)
	if err != nil {
		h.log.ErrorContext(ctx, "failed to list active draw", "error", err)
		helpers.ErrorMessage(w, "list error", http.StatusBadRequest, err)

		return
	}

	helpers.SuccessMessage(w, "draws", list)
}
