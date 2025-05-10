package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (h *handler) ListActiveDraws(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.log.InfoContext(ctx, "list active draws", "id", r.PathValue("draw_id"))

	// Получаем список активных тиражей
	list, err := h.service.ListActiveDraws(ctx)
	if err != nil {
		h.log.ErrorContext(ctx, "failed to list active draw", "error", err)
		http.Error(w, fmt.Sprintf("error on list: %s", err.Error()), http.StatusBadRequest)

		return
	}

	// В случае успеха, подготавливаем ответ
	result, err := json.Marshal(list)
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
