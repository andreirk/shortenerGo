package stat

import (
	"fmt"
	"go/adv-demo/config"
	"go/adv-demo/pkg/response"
	"net/http"
	"time"
)

const (
	GroupByMonth = "month"
	GroupByDay   = "day"
)

type StatHandlerDeps struct {
	StatRepository *StatRepository
	Config         *config.Config
}

type StatHandler struct {
	StatRepository *StatRepository
}

func NewStatHandler(router *http.ServeMux, deps StatHandlerDeps) *StatHandler {
	handler := &StatHandler{
		StatRepository: deps.StatRepository,
	}

	router.HandleFunc("GET /stat", handler.GetStat())

	return handler
}

func (handler *StatHandler) GetStat() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		from, err := time.Parse("2006-01-02", r.URL.Query().Get("from"))
		if err != nil {
			http.Error(w, "Invalid from param", http.StatusBadRequest)
			return
		}
		to, err := time.Parse("2006-01-02", r.URL.Query().Get("to"))
		if err != nil {
			http.Error(w, "Invalid to param", http.StatusBadRequest)
			return
		}
		by := r.URL.Query().Get("by")
		if by != GroupByDay && by != GroupByMonth {
			http.Error(w, "Invalid by param", http.StatusBadRequest)
			return
		}
		fmt.Printf("from: %v, to: %v\n, by: %s", from, to, by)
		stats := handler.StatRepository.GetStats(by, from, to)
		response.JsonResponse(w, stats, http.StatusOK)
	}
}
