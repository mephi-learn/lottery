package controller

import (
	"encoding/json"
	"fmt"
	"homework/internal/auth"
	"net/http"
)

func (h *handler) signUp(w http.ResponseWriter, r *http.Request) {
	var user auth.SignUpInput

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		NewErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	userId, err := h.service.SignUp(r.Context(), &user)
	if err != nil {
		NewErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("user was created, id = %d", userId)))
}
