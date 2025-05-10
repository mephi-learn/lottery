package controller

import (
	"encoding/json"
	"fmt"
	"homework/internal/models"
	"net/http"
	"strconv"
)

type getTicketByIdResponse struct {
	responseTicket
}

func (h *handler) GetTicketById(w http.ResponseWriter, r *http.Request) {
	_, err := models.UserFromContext(r.Context())
	if err != nil {
		http.Error(w, "authenticate need", http.StatusBadRequest)
		return
	}

	ticketId, err := strconv.Atoi(r.PathValue("ticket_id"))
	if err != nil {
		http.Error(w, fmt.Sprintf("invalid draw: %s", r.PathValue("ticket_id")), http.StatusBadRequest)
		return
	}

	ticket, err := h.service.GetTicketById(r.Context(), ticketId)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed get ticket: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	ticketCombination, err := models.ParseTicketCombination(ticket.Data)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed get ticket combination: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	response := getTicketByIdResponse{
		responseTicket: responseTicket{
			Id:          ticket.Id,
			StatusName:  ticket.Status.String(),
			DrawId:      ticket.DrawId,
			UserId:      ticket.UserId,
			Combination: ticketCombination,
			Cost:        ticket.Cost,
		},
	}

	out, err := json.Marshal(response)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed create response: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(out))
}
