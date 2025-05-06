package controller

import (
	"encoding/json"
	"homework/internal/models"
	"net/http"
)

func (h *handler) signIn(w http.ResponseWriter, r *http.Request) {
	var signIn models.SignInInput

	err := json.NewDecoder(r.Body).Decode(&signIn)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	signedToken, err := h.service.SignIn(r.Context(), &signIn)
	if err != nil {
		http.Error(w, "user was not found", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`
{
	"token": "` + signedToken + `"
}`))
}
