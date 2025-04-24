package controller

import (
	"encoding/json"
	"fmt"
	"homework/internal/models"
	"net/http"
)

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

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(fmt.Sprintf("draw was created, id = %d", drawId)))
}
