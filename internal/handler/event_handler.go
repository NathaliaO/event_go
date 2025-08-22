package handler

import (
	"encoding/json"
	"net/http"

	"github.com/nathaliaoliveira/goapp/internal/domain"
	"github.com/nathaliaoliveira/goapp/internal/service"
)

type EventHandler struct {
    eventService service.EventService
}

func NewEventHandler(eventService service.EventService) *EventHandler {
    return &EventHandler{
        eventService: eventService,
    }
}

func (h *EventHandler) CreateEvents(w http.ResponseWriter, r *http.Request) {
	var eventsReq domain.EventsRequest
	if err := json.NewDecoder(r.Body).Decode(&eventsReq); err != nil {
		http.Error(w, "Dados inválidos", http.StatusBadRequest)
		return
	}
	
	result, err := h.eventService.ProcessEvents(eventsReq.Events)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(result)
}

func (h *EventHandler) GetDailyStats(w http.ResponseWriter, r *http.Request) {
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")
	site := r.URL.Query().Get("site")
	
	stats, err := h.eventService.GetDailyStats(startDate, endDate, site)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	
	response := domain.Response{
		Message: "Estatísticas diárias por site",
		Data:    stats,
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *EventHandler) handleServiceError(w http.ResponseWriter, err error) {
    switch e := err.(type) {
    case *service.ValidationError:
        http.Error(w, e.Error(), http.StatusBadRequest)
    default:
        http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
    }
} 