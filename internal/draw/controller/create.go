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
	user, err := models.UserFromContext(r.Context())
	if err != nil {
		http.Error(w, "authenticate need", http.StatusBadRequest)
		return
	}

	if !user.Admin {
		http.Error(w, "permission denied, admin only area", http.StatusForbidden)
		return
	}

	// Парсим входные данные
	draw := models.DrawInput{}
	err = json.NewDecoder(r.Body).Decode(&draw)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed decore request: %s", err.Error()), http.StatusBadRequest)
		return
	}

	drawId, err := h.service.CreateDraw(r.Context(), &draw)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed create draw: %s", err.Error()), http.StatusBadRequest)
		return
	}

	resp := drawResponse{
		Message: "draw has been created",
		DrawId:  drawId,
	}

	result, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed create response: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(result)
}
