package controller

import (
	"fmt"
	"homework/internal/helpers"
	"homework/internal/models"
	"net/http"
	"strconv"
)

func (h *handler) GetTicketById(w http.ResponseWriter, r *http.Request) {
	ticketId, err := strconv.Atoi(r.PathValue("ticket_id"))
	if err != nil {
		helpers.ErrorMessage(w, fmt.Sprintf("invalid ticket id: %s", r.PathValue("ticket_id")), http.StatusBadRequest, err)
		return
	}

	ticket, err := h.service.GetTicketById(r.Context(), ticketId)
	if err != nil {
		helpers.ErrorMessage(w, "failed to get ticket", http.StatusBadRequest, err)
		return
	}

	ticketCombination, err := models.ParseTicketCombination(ticket.Data)
	if err != nil {
		helpers.ErrorMessage(w, "failed get ticket combination", http.StatusBadRequest, err)
		return
	}

	response := responseTicket{
		Id:          ticket.Id,
		StatusName:  ticket.Status.String(),
		DrawId:      ticket.DrawId,
		UserId:      ticket.UserId,
		Combination: ticketCombination,
		Cost:        ticket.Cost,
	}

	helpers.SuccessMessage(w, "ticket", response)
}
