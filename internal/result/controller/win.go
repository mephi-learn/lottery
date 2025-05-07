package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func (h *handler) GetDrawWinResults(w http.ResponseWriter, r *http.Request) {
	// Парсим входные данные
	drawId, err := strconv.Atoi(r.PathValue("draw_id"))
	if err != nil {
		http.Error(w, fmt.Sprintf("invalid id: %s", r.PathValue("draw_id")), http.StatusBadRequest)
		return
	}

	list, err := h.service.GetDrawWinResults(r.Context(), drawId)
	if err != nil {
		h.log.Error("failed to get draw results", "err", err)
		http.Error(w, fmt.Sprintf("failed to get draw results: %s", err.Error()), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	//comb, err := json.Marshal(combination)
	//if err != nil {
	//	http.Error(w, fmt.Sprintf("failed create response: %s", err.Error()), http.StatusInternalServerError)
	//	return
	//}
	//_, _ = w.Write(comb)
	//_, _ = w.Write([]byte{'\n', '\n'})

	stats, err := json.Marshal(list.Statistic)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed create response: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(stats)
	_, _ = w.Write([]byte{'\n', '\n'})

	out, err := json.MarshalIndent(list.WinTickets, "", "  ")
	if err != nil {
		http.Error(w, fmt.Sprintf("failed create response: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	_, err = w.Write(out)
}
