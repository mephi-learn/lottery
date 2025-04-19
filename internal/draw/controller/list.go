package controller

import "net/http"

func (h *handler) ListActiveDraw(writer http.ResponseWriter, request *http.Request) {
	h.service.ListActiveDraw(request.Context())
}
