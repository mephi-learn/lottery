package controller

import (
	"net/http"
	"time"
)

func (h *handler) RegisterInvoice(writer http.ResponseWriter, request *http.Request) {
	h.service.RegisterInvoice(request.Context(), time.Now())
}
