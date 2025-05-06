package controller

import (
	"encoding/json"
	"net/http"
)

func (h *handler) ExportDraws(w http.ResponseWriter, r *http.Request) {
	draws, err := h.service.ExportDraws(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(draws)
	if err != nil {
		h.log.Error("failed to marshal response", "err", err)
		http.Error(w, "failed to marshal response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(data)
}
