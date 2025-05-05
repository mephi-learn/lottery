package controller

import (
	"encoding/json"
	"fmt"
	"homework/internal/models"
	"net/http"
	"strconv"
)

type listAvailableTicketsResponse struct {
	Tickets []models.Ticket `json:"tickets"`
}

func (h *handler) ListAvailableTickets(w http.ResponseWriter, r *http.Request) {
	drawId, err := strconv.Atoi(r.PathValue("drawId"))
	if err != nil {
		http.Error(w, fmt.Sprintf("invalid draw: %s", r.PathValue("drawId")), http.StatusBadRequest)
		return
	}

	tickets, err := h.service.ListAvailableTicketsByDrawId(r.Context(), drawId)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed get ticket: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	response := listAvailableTicketsResponse{
		Tickets: tickets,
	}

	out, err := json.Marshal(response)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed create response: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(out))
}
