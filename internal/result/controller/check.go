package controller

import (
	"fmt"
	"homework/internal/helpers"
	"homework/internal/models"
	"net/http"
	"strconv"
)

func (h *handler) CheckTicketResult(w http.ResponseWriter, r *http.Request) {
	user, err := models.UserFromContext(r.Context())
	if err != nil {
		helpers.ErrorMessage(w, "authenticate need", http.StatusBadRequest, nil)
		return
	}

	// Парсим входные данные
	ticketId, err := strconv.Atoi(r.PathValue("ticket_id"))
	if err != nil {
		helpers.ErrorMessage(w, fmt.Sprintf("no such ticket for the given user: %s", r.PathValue("ticket_id")), http.StatusBadRequest, err)
		return
	}

	result, err := h.service.CheckTicketResult(r.Context(), ticketId, user.ID)
	if err != nil {
		h.log.Error("failed to check ticket result", "err", err)
		helpers.ErrorMessage(w, "failed to check ticket result", http.StatusBadRequest, err)

		return
	}

	helpers.SuccessMessage(w, "result", result)
}

func (h *handler) CheckTicketsResult(w http.ResponseWriter, r *http.Request) {
	user, err := models.UserFromContext(r.Context())
	if err != nil {
		helpers.ErrorMessage(w, "authenticate need", http.StatusBadRequest, nil)
		return
	}

	result, err := h.service.CheckTicketsResult(r.Context(), user.ID)
	if err != nil {
		h.log.Error("failed to check user tickets", "err", err)
		helpers.ErrorMessage(w, "failed to check user tickets result", http.StatusBadRequest, err)

		return
	}

	helpers.SuccessMessage(w, "result", result)
}
