package controller

import (
	"fmt"
	"homework/internal/helpers"
	"homework/internal/models"
	"net/http"
	"strconv"
)

type responseTicket struct {
	Id          int     `json:"id"`
	StatusName  string  `json:"status_name"`
	DrawId      int     `json:"draw_id"`
	UserId      int     `json:"user_id"`
	Combination []int   `json:"combination"`
	Cost        float64 `json:"cost"`
}

type listAvailableTicketsResponse struct {
	Tickets []responseTicket `json:"tickets"`
}

func (h *handler) ListAvailableTickets(w http.ResponseWriter, r *http.Request) {
	drawId, err := strconv.Atoi(r.PathValue("draw_id"))
	if err != nil {
		helpers.ErrorMessage(w, fmt.Sprintf("invalid draw id: %s", r.PathValue("draw_id")), http.StatusBadRequest, err)
		return
	}

	tickets, err := h.service.ListAvailableTicketsByDrawId(r.Context(), drawId)
	if err != nil {
		helpers.ErrorMessage(w, "failed to get ticket", http.StatusBadRequest, err)
		return
	}

	ticketsVal := make([]responseTicket, 0, len(tickets))
	for _, t := range tickets {
		if t != nil {
			ticketCombination, err := models.ParseTicketCombination(t.Data)
			if err != nil {
				helpers.ErrorMessage(w, "failed get ticket combination", http.StatusBadRequest, err)
				return
			}

			ticketsVal = append(ticketsVal, responseTicket{
				Id:          t.Id,
				StatusName:  t.Status.String(),
				DrawId:      t.DrawId,
				UserId:      t.UserId,
				Combination: ticketCombination,
				Cost:        t.Cost,
			})
		}
	}

	helpers.SuccessMessage(w, "tickets", ticketsVal)
}
