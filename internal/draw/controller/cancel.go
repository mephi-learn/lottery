package controller

import (
	"encoding/json"
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
		h.log.Error("failed to cancel draw", "error", err)
		http.Error(w, fmt.Sprintf("failed to cancel draw: %s", err.Error()), http.StatusInternalServerError)

		return
	}

	resp := drawResponse{
		Message: "draw has been canceled",
		DrawId:  id,
	}

	result, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed create response: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(result)
}
