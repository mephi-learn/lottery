package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (h *handler) ListActiveDraw(w http.ResponseWriter, r *http.Request) {
	list, err := h.service.ListActiveDraw(r.Context())
	if err != nil {
		http.Error(w, fmt.Sprintf("error on list: %s", err.Error()), http.StatusBadRequest)
		return
	}

	out, err := json.Marshal(list)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed create response: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(out)
}
