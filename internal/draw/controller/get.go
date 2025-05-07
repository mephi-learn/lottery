package controller

import (
	"encoding/json"
	"fmt"
	"homework/internal/models"
	"net/http"
	"strconv"
)

func (h *handler) GetDraw(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.log.InfoContext(ctx, "get draw", "id", r.PathValue("draw_id"))

	// Получаем идентификатор тиража
	drawId, err := strconv.Atoi(r.PathValue("draw_id"))
	if err != nil {
		h.log.ErrorContext(ctx, "invalid draw id", "id", r.PathValue("draw_id"), "error", err)
		http.Error(w, fmt.Sprintf("invalid id: %s", r.PathValue("draw_id")), http.StatusBadRequest)

		return
	}

	// Получаем информацию по тиражу
	draw, err := h.service.GetDraw(ctx, drawId)
	if err != nil {
		h.log.ErrorContext(ctx, "failed to get draw", "error", err)
		http.Error(w, fmt.Sprintf("error on get draw: %s", err.Error()), http.StatusBadRequest)

		return
	}

	// На основе тиража получаем лотерею
	lottery, err := h.service.LotteryByType(draw.LotteryType)
	if err != nil {
		h.log.ErrorContext(ctx, "failed to detect lottery", "error", err)
		http.Error(w, fmt.Sprintf("lottery unknown type: %s", err.Error()), http.StatusBadRequest)

		return
	}

	// В случае успеха, подготавливаем ответ
	status := models.DrawStatus(draw.StatusId)
	drawOut := models.DrawOutput{
		Id:        draw.Id,
		Status:    status.String(),
		Lottery:   lottery.Type(),
		SaleDate:  draw.SaleDate,
		StartDate: draw.StartDate,
	}
	result, err := json.Marshal(drawOut)
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
