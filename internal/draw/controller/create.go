package controller

import (
	"net/http"
	"time"
)

func (h *handler) CreateDraw(writer http.ResponseWriter, request *http.Request) {
	h.service.CreateDraw(request.Context(), time.Now(), time.Now().Add(time.Minute))
}
