package controller

import (
	"fmt"
	"homework/internal/models"
	"net/http"
	"strconv"
)

func (h *handler) CreateTickets(w http.ResponseWriter, r *http.Request) {
	user, err := models.UserFromContext(r.Context())
	if err != nil {
		http.Error(w, "authenticate need", http.StatusBadRequest)
		return
	}

	if !user.Admin {
		http.Error(w, "permission denied, admin only area", http.StatusForbidden)
		return
	}

	drawId, err := strconv.Atoi(r.PathValue("drawId"))
	if err != nil {
		http.Error(w, fmt.Sprintf("invalid draw: %s", r.PathValue("drawId")), http.StatusBadRequest)
		return
	}

	num, err := strconv.Atoi(r.PathValue("num"))
	if err != nil {
		http.Error(w, fmt.Sprintf("invalid num: %s", r.PathValue("num")), http.StatusBadRequest)
		return
	}

	list, err := h.service.CreateTickets(r.Context(), drawId, num)

	// TODO: вернуть число созданных билетов
	out := fmt.Sprintf("created %d tickets", len(list))

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(out))
}
