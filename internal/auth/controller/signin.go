package controller

import (
	"encoding/json"
	"homework/internal/auth"
	"net/http"
)

func (h *handler) signIn(w http.ResponseWriter, r *http.Request) {
	var signIn auth.SignInInput

	err := json.NewDecoder(r.Body).Decode(&signIn)
	if err != nil {
		NewErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	signedToken, err := h.service.SignIn(r.Context(), &signIn)
	if err != nil {
		NewErrorResponse(w, http.StatusBadRequest, "user was not found")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(signedToken))
}
