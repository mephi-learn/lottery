package controller

import (
	"fmt"
	"homework/internal/helpers"
	"net/http"
	"strconv"
)

func (h *handler) CancelDraw(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.log.InfoContext(ctx, "cancel draw", "id", r.PathValue("draw_id"))

	// Получаем идентификатор тиража
	drawId, err := strconv.Atoi(r.PathValue("draw_id"))
	if err != nil {
		h.log.ErrorContext(ctx, "invalid draw id", "id", r.PathValue("draw_id"), "error", err)
		helpers.ErrorMessage(w, fmt.Sprintf("invalid id: %s", r.PathValue("draw_id")), http.StatusBadRequest, err)

		return
	}

	// Отменяем тираж
	if err := h.service.CancelDraw(ctx, drawId); err != nil {
		h.log.ErrorContext(ctx, "failed to cancel draw", "error", err)
		helpers.ErrorMessage(w, "failed to cancel draw", http.StatusInternalServerError, err)

		return
	}

	// Выдаём ответ
	helpers.SuccessMessage(w, "draw has been canceled", map[string]any{"draw id": drawId})
}
