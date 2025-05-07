package controller

import (
	"encoding/json"
	"fmt"
	"homework/internal/models"
	"net/http"
)

type drawResponse struct {
	Message string `json:"message"`
	DrawId  int    `json:"draw_id"`
}

func (h *handler) CreateDraw(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.log.InfoContext(ctx, "create draw", "id", r.PathValue("draw_id"))

	// Получаем идентификатор тиража
	draw := models.DrawInput{}
	if err := json.NewDecoder(r.Body).Decode(&draw); err != nil {
		h.log.ErrorContext(ctx, "invalid draw id", "id", r.PathValue("draw_id"), "error", err)
		http.Error(w, fmt.Sprintf("failed decore request: %s", err.Error()), http.StatusBadRequest)

		return
	}

	// Создаём тираж
	drawId, err := h.service.CreateDraw(ctx, &draw)
	if err != nil {
		h.log.ErrorContext(ctx, "failed to create draw", "error", err)
		http.Error(w, fmt.Sprintf("failed create draw: %s", err.Error()), http.StatusBadRequest)

		return
	}

	// В случае успеха, подготавливаем ответ
	resp := drawResponse{
		Message: "draw has been created",
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
