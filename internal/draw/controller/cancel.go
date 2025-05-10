package controller

import (
	"encoding/json"
	"fmt"
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
		http.Error(w, fmt.Sprintf("invalid id: %s", r.PathValue("draw_id")), http.StatusBadRequest)

		return
	}

	// Отменяем тираж
	if err := h.service.CancelDraw(ctx, drawId); err != nil {
		h.log.ErrorContext(ctx, "failed to cancel draw", "error", err)
		http.Error(w, fmt.Sprintf("failed to cancel draw: %s", err.Error()), http.StatusInternalServerError)

		return
	}

	// В случае успеха, подготавливаем ответ
	resp := drawResponse{
		Message: "draw has been canceled",
		DrawId:  drawId,
	}
	result, err := json.Marshal(resp)
	if err != nil {
		h.log.ErrorContext(ctx, "failed to encode json response", "error", err)
		http.Error(w, fmt.Sprintf("failed create response: %s", err.Error()), http.StatusInternalServerError)

		return
	}

	// И возвращаем его
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(result)
}
