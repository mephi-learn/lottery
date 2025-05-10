package controller

import (
	"encoding/json"
	"fmt"
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
		http.Error(w, fmt.Sprintf("invalid draw: %s", r.PathValue("draw_id")), http.StatusBadRequest)
		return
	}

	tickets, err := h.service.ListAvailableTicketsByDrawId(r.Context(), drawId)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed get ticket: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	ticketsVal := make([]responseTicket, 0, len(tickets))
	for _, t := range tickets {
		if t != nil {
			ticketCombination, err := models.ParseTicketCombination(t.Data)
			if err != nil {
				http.Error(w, fmt.Sprintf("failed get ticket combination: %s", err.Error()), http.StatusInternalServerError)
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
