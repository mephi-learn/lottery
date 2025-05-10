package controller

import (
	"encoding/json"
	"homework/internal/helpers"
	"homework/internal/models"
	"net/http"
)

func (h *handler) signIn(w http.ResponseWriter, r *http.Request) {
	var signIn models.SignInInput

	err := json.NewDecoder(r.Body).Decode(&signIn)
	if err != nil {
		helpers.ErrorMessage(w, "invalid json data", http.StatusBadRequest, err)
		return
	}

	signedToken, err := h.service.SignIn(r.Context(), &signIn)
	if err != nil {
		helpers.ErrorMessage(w, "user was not found", http.StatusBadRequest, nil)
		return
	}

	helpers.SuccessMessage(w, "token create", map[string]any{"token": signedToken})
}
