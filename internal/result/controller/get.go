package controller

import (
	"fmt"
	"net/http"
	"strconv"
)

// type cancelDraw struct {
// 	Id int `json:"id"`
// }

func (h *handler) GetDrawResults(w http.ResponseWriter, r *http.Request) {
	// user, err := models.UserFromContext(r.Context())
	// if err != nil {
	// 	http.Error(w, "authenticate need", http.StatusBadRequest)
	// 	return
	// }

	// if !user.Admin {
	// 	http.Error(w, "permission denied, admin only area", http.StatusForbidden)
	// 	return
	// }

	// Парсим входные данные
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, fmt.Sprintf("invalid id: %s", r.PathValue("id")), http.StatusBadRequest)
		return
	}

	result, err := h.service.GetDrawResults(r.Context(), id); 
	if err != nil {
		h.log.Error("failed to get draw results", "err", err)
		http.Error(w, fmt.Sprintf("failed to get draw results: %w", err.Error()), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(fmt.Sprintf("draw results: %d", result)))
}
