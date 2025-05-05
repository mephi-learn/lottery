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
	drawId, err := strconv.Atoi(r.PathValue("draw_id"))
	if err != nil {
		http.Error(w, fmt.Sprintf("invalid draw: %s", r.PathValue("draw_id")), http.StatusBadRequest)
		return
	}

	tickets, err := h.service.ListAvailableTicketsByDrawId(r.Context(), drawId)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed get ticket: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	ticketsVal := make([]models.Ticket, 0, len(tickets))
	for _, t := range tickets {
		if t != nil {
			ticketsVal = append(ticketsVal, *t)
		}
	}

	response := listAvailableTicketsResponse{
		Tickets: ticketsVal,
	}

	out, err := json.Marshal(response)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed create response: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(out))
}
