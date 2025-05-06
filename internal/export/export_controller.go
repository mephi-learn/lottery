package export

import (
	"encoding/csv"
	"fmt"
	"homework/internal/auth"
	"homework/internal/export/service"
	"homework/pkg/errors"
	"homework/pkg/log"
	"net/http"
	"strconv"
	"strings"
)

type handler struct {
	log     log.Logger
	service *service.Service
}

type HandlerOption func(*handler)

func NewHandler(opts ...HandlerOption) (*handler, error) {
	h := &handler{}

	for _, opt := range opts {
		opt(h)
	}

	if h.log == nil {
		return nil, errors.New("logger is missing")
	}

	if h.service == nil {
		return nil, errors.New("service is missing")
	}

	return h, nil
}

func WithLogger(logger log.Logger) HandlerOption {
	return func(o *handler) {
		o.log = logger
	}
}

func WithService(svc *service.Service) HandlerOption {
	return func(o *handler) {
		o.service = svc
	}
}

func (h *handler) WithRouter(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/admin/export/draws", auth.Authenticated(h.ExportDraws))
}

func (h *handler) ExportDraws(w http.ResponseWriter, r *http.Request) {
	draws, err := h.service.ExportDraws(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/csv")
	writer := csv.NewWriter(w)
	defer writer.Flush()

	// Write header
	writer.Write([]string{"draw_id", "lottery_type", "winning_combination", "winner_count"})

	// Write data
	for _, draw := range draws {
		writer.Write([]string{
			strconv.FormatInt(draw.DrawID, 10),
			draw.LotteryType,
			strings.Trim(strings.Join(strings.Fields(fmt.Sprint(draw.WinningCombination)), " "), "[]"),
			strconv.Itoa(draw.WinnerCount),
		})
	}
}
