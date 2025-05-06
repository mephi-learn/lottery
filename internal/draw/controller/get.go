package controller

import (
	"encoding/json"
	"fmt"
	"homework/internal/models"
	"net/http"
	"strconv"
)

func (h *handler) GetDraw(w http.ResponseWriter, r *http.Request) {
	// Парсим входные данные
	drawId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, fmt.Sprintf("invalid id: %s", r.PathValue("id")), http.StatusBadRequest)
		return
	}

	draw, err := h.service.GetDraw(r.Context(), drawId)
	if err != nil {
		http.Error(w, fmt.Sprintf("error on get draw: %s", err.Error()), http.StatusBadRequest)
		return
	}

	lottery, err := h.service.LotteryByType(draw.LotteryType)
	if err != nil {
		http.Error(w, fmt.Sprintf("lottery unknown type: %s", err.Error()), http.StatusBadRequest)
		return
	}
	status := models.DrawStatus(draw.StatusId)

	drawOut := models.DrawOutput{
		Id:        draw.Id,
		Status:    status.String(),
		Lottery:   lottery.Type(),
		SaleDate:  draw.SaleDate,
		StartDate: draw.StartDate,
	}
	out, err := json.Marshal(drawOut)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed create response: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(out)
}
