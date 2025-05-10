package controller

import (
	"encoding/json"
	"fmt"
	"homework/internal/models"
	"net/http"
	"strconv"
)

func (h *handler) CheckTicketResult(w http.ResponseWriter, r *http.Request) {
	user, err := models.UserFromContext(r.Context())
	if err != nil {
		http.Error(w, "authentication needed", http.StatusBadRequest)
		return
	}

	// Парсим входные данные
	ticketId, err := strconv.Atoi(r.PathValue("ticket_id"))
	if err != nil {
		http.Error(w, fmt.Sprintf("no such ticket for the given user: %s", r.PathValue("ticket_id")), http.StatusBadRequest)
		return
	}

	result, err := h.service.CheckTicketResult(r.Context(), ticketId, user.ID)
	if err != nil {
		h.log.Error("failed to check ticket result", "err", err)
		http.Error(w, fmt.Sprintf("failed to check ticket result: %s", err.Error()), http.StatusInternalServerError)

		return
	}

	slice := []*models.TicketResult{result}

	responsePayload := map[string]interface{}{
		"result": slice,
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

func (h *handler) CheckTicketsResult(w http.ResponseWriter, r *http.Request) {
	user, err := models.UserFromContext(r.Context())
	if err != nil {
		http.Error(w, "authentication needed", http.StatusBadRequest)
		return
	}

	result, err := h.service.CheckTicketsResult(r.Context(), user.ID)
	if err != nil {
		h.log.Error("failed to check user tickets", "err", err)
		http.Error(w, fmt.Sprintf("failed to check user tickets result: %s", err.Error()), http.StatusInternalServerError)

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
