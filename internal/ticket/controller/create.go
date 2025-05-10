package controller

import (
	"fmt"
	"homework/internal/helpers"
	"net/http"
	"strconv"
)

func (h *handler) CreateTickets(w http.ResponseWriter, r *http.Request) {
	drawId, err := strconv.Atoi(r.PathValue("draw_id"))
	if err != nil {
		helpers.ErrorMessage(w, fmt.Sprintf("invalid draw id: %s", r.PathValue("draw_id")), http.StatusBadRequest, err)
		return
	}

	num, err := strconv.Atoi(r.PathValue("count"))
	if err != nil {
		helpers.ErrorMessage(w, fmt.Sprintf("invalid count: %s", r.PathValue("count")), http.StatusBadRequest, err)
		return
	}

	list, err := h.service.CreateTickets(r.Context(), drawId, num)
	if err != nil {
		helpers.ErrorMessage(w, "failed create tickets", http.StatusBadRequest, err)
		return
	}

	helpers.SuccessMessageWithCode(w, "tickets created", list, http.StatusCreated)
}
