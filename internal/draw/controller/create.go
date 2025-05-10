package controller

import (
	"encoding/json"
	"homework/internal/helpers"
	"homework/internal/models"
	"net/http"
)

func (h *handler) CreateDraw(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.log.InfoContext(ctx, "create draw", "id", r.PathValue("draw_id"))

	// Получаем идентификатор тиража
	draw := models.DrawInput{}
	if err := json.NewDecoder(r.Body).Decode(&draw); err != nil {
		h.log.ErrorContext(ctx, "invalid draw id", "id", r.PathValue("draw_id"), "error", err)
		helpers.ErrorMessage(w, "invalid json data", http.StatusBadRequest, err)

		return
	}

	if draw.Cost <= 0 {
		helpers.ErrorMessage(w, "incorrect cost", http.StatusBadRequest, nil)
		return
	}

	// Создаём тираж
	drawId, err := h.service.CreateDraw(ctx, &draw)
	if err != nil {
		h.log.ErrorContext(ctx, "failed to create draw", "error", err)
		helpers.ErrorMessage(w, "failed create draw", http.StatusBadRequest, err)

		return
	}

	// Выдаём ответ
	helpers.SuccessMessageWithCode(w, "draw has been created", map[string]any{"draw id": drawId}, http.StatusCreated)
}
