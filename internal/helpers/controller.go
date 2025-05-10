package helpers

import (
	"encoding/json"
	"homework/internal/models"
	"net/http"
)

func ErrorMessage(w http.ResponseWriter, message string, code int, err error) {
	// Формируем сообщение
	if err != nil {
		message += ": " + err.Error()
	}
	resp := models.ResultErrorMessage{
		Error: message,
	}
	result, _ := json.Marshal(resp)

	http.Error(w, string(result), code)
}

func SuccessMessage(w http.ResponseWriter, message string, data any) {
	SuccessMessageWithCode(w, message, data, http.StatusOK)
}

func SuccessMessageWithCode(w http.ResponseWriter, message string, data any, code int) {
	// Формируем сообщение
	resp := models.SuccessMessage{
		Message: message,
		Data:    data,
	}
	result, _ := json.Marshal(resp)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_, _ = w.Write(result)
}
