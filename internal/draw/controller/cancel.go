package controller

import (
	"fmt"
	"homework/internal/models"
	"net/http"
	"strconv"
)

type cancelDraw struct {
	Id int `json:"id"`
}

func (h *handler) CancelDraw(w http.ResponseWriter, r *http.Request) {
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
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, fmt.Sprintf("invalid id: %s", r.PathValue("id")), http.StatusBadRequest)
		return
	}

	if err := h.service.CancelDraw(r.Context(), id); err != nil {
		h.log.Error("failed to cancel draw", "err", err)
		http.Error(w, fmt.Sprintf("failed to cancel draw: %w", err.Error()), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(fmt.Sprintf("draw was canceled, id = %d", id)))
}
