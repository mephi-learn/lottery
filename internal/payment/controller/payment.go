package controller

import (
	"net/http"
)

func (h *handler) RegisterPayment(writer http.ResponseWriter, request *http.Request) {
	mockPaymentData := 1000.5 // просто временный мок
	h.service.RegisterPayment(request.Context(), mockPaymentData)
}
