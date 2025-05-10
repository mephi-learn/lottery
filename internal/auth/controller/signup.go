package controller

import (
	"encoding/json"
	"homework/internal/helpers"
	"homework/internal/models"
	"net/http"
)

func (h *handler) signUp(w http.ResponseWriter, r *http.Request) {
	var user models.SignUpInput

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		helpers.ErrorMessage(w, "invalid json data", http.StatusBadRequest, err)
		return
	}

	userId, err := h.service.SignUp(r.Context(), &user)
	if err != nil {
		helpers.ErrorMessage(w, "sign up error", http.StatusBadRequest, nil)
		return
	}

	helpers.SuccessMessage(w, "user was created", map[string]any{"user id": userId})
}
