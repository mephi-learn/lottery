package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func (h *handler) GetDrawResults(w http.ResponseWriter, r *http.Request) {
	// Парсим входные данные
	id, err := strconv.Atoi(r.PathValue("draw_id"))
	if err != nil {
		http.Error(w, fmt.Sprintf("invalid id: %s", r.PathValue("draw_id")), http.StatusBadRequest)
		return
	}

	result, err := h.service.GetDrawResults(r.Context(), id)
	if err != nil {
		h.log.Error("failed to get draw results", "err", err)
		http.Error(w, fmt.Sprintf("failed to get draw results: %s", err.Error()), http.StatusInternalServerError)

		return
	}

	responsePayload := map[string]interface{}{
		"result": result,
	}
	data, err := json.Marshal(responsePayload)
	if err != nil {
		h.log.Error("failed to marshal response", "err", err)
		http.Error(w, "failed to marshal response", http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(data)
}
