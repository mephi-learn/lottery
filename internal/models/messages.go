package models

type ResultErrorMessage struct {
	Error string `json:"error"`
}

type SuccessMessage struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}
