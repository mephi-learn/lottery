package controller

import (
	"net/http"
)

func (h *handler) CancelDraw(writer http.ResponseWriter, request *http.Request) {
	err := h.service.CancelDraw(request.Context(), 0)
	if err != nil {
		h.log.Error("failed to cancel draw", "err", err)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}
