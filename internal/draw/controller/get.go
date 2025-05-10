package controller

import (
	"fmt"
	"homework/internal/helpers"
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
		helpers.ErrorMessage(w, fmt.Sprintf("invalid draw id: %s", r.PathValue("draw_id")), http.StatusBadRequest, err)

		return
	}

	// Получаем информацию по тиражу
	draw, err := h.service.GetDraw(ctx, drawId)
	if err != nil {
		h.log.ErrorContext(ctx, "failed to get draw", "error", err)
		helpers.ErrorMessage(w, "failed to get draw", http.StatusBadRequest, err)

		return
	}

	// На основе тиража получаем лотерею
	lottery, err := h.service.LotteryByType(draw.LotteryType)
	if err != nil {
		h.log.ErrorContext(ctx, "failed to detect lottery", "error", err)
		helpers.ErrorMessage(w, "lottery unknown type", http.StatusBadRequest, err)

		return
	}

	// В случае успеха, подготавливаем и выдаём ответ
	status := models.DrawStatus(draw.StatusId)
	drawOut := models.DrawOutput{
		Id:        draw.Id,
		Status:    status.String(),
		Lottery:   lottery.Type(),
		SaleDate:  draw.SaleDate,
		StartDate: draw.StartDate,
	}

	helpers.SuccessMessage(w, "draw", drawOut)
}
