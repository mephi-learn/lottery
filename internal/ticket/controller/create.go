package controller

import (
	"encoding/json"
	"net/http"
	"strconv"
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

	drawID, err := strconv.ParseInt(input.DrawID, 10, 64)
	if err != nil {
		http.Error(w, "Invalid draw_id", http.StatusBadRequest)
		return
	}

	// TODO: get userID from context
	userID := int64(1)

	ticket, err := h.service.CreateTicket(r.Context(), userID, drawID, input.Numbers, false)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"ticket_id": ticket.ID})
}
