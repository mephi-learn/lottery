package controller

import (
	"encoding/json"
	"homework/internal/models"
	"net/http"
)

type ticketInput struct {
	DrawID  string `json:"draw_id"`
	Numbers []int  `json:"numbers"`
}

func (h *handler) CreateTicket(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var input ticketInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ticket := &models.TicketRequest{
		DrawID:  input.DrawID,
		Numbers: input.Numbers,
	}

	ticketId, err := h.service.CreateTicket(r.Context(), ticket)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"ticket_id": ticketId})
}
