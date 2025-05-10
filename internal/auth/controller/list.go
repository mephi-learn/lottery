package controller

import (
	"homework/internal/helpers"
	"net/http"
)

func (h *handler) list(w http.ResponseWriter, r *http.Request) {
	signedToken, err := h.service.List(r.Context())
	if err != nil {
		helpers.ErrorMessage(w, "user was not found", http.StatusBadRequest, nil)
		return
	}

	helpers.SuccessMessage(w, "token create", signedToken)
}
