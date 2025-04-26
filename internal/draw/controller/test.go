package controller

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strconv"
)

func (h *handler) Drawing(w http.ResponseWriter, r *http.Request) {
	drawId, err := strconv.Atoi(r.PathValue("draw"))
	if err != nil {
		http.Error(w, fmt.Sprintf("invalid draw: %s", r.PathValue("draw")), http.StatusBadRequest)
		return
	}

	combination := make([]int, 5)
	maxDigitIndex := big.NewInt(35 - 1)
	for i := range 5 {
		randomNumber, _ := rand.Int(rand.Reader, maxDigitIndex)
		combination[i] = int(randomNumber.Int64()) + 1
	}

	list, err := h.service.Drawing(r.Context(), drawId, combination)
	if err != nil {
		http.Error(w, fmt.Sprintf("error: %s", err.Error()), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)

	comb, err := json.Marshal(combination)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed create response: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(comb)
	_, _ = w.Write([]byte{'\n', '\n'})

	stat := map[string]int{}
	for key, values := range list {
		stat[key] = len(values)
	}
	stats, err := json.Marshal(stat)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed create response: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(stats)
	_, _ = w.Write([]byte{'\n', '\n'})

	out, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		http.Error(w, fmt.Sprintf("failed create response: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(out)
}

func (h *handler) Generate(w http.ResponseWriter, r *http.Request) {
	drawId, err := strconv.Atoi(r.PathValue("draw"))
	if err != nil {
		http.Error(w, fmt.Sprintf("invalid draw: %s", r.PathValue("draw")), http.StatusBadRequest)
		return
	}

	num, err := strconv.Atoi(r.PathValue("num"))
	if err != nil {
		http.Error(w, fmt.Sprintf("invalid num: %s", r.PathValue("num")), http.StatusBadRequest)
		return
	}

	list, err := h.service.CreateTickets(r.Context(), drawId, num)

	out, err := json.Marshal(list)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed create response: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(out)
}
