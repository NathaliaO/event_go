package handler

import (
    "encoding/json"
    "log"
    "net/http"

    "github.com/nathaliaoliveira/goapp/internal/domain"
    "github.com/nathaliaoliveira/goapp/internal/service"
)

type HealthHandler struct {
    healthService service.HealthService
}

func NewHealthHandler(healthService service.HealthService) *HealthHandler {
    return &HealthHandler{
        healthService: healthService,
    }
}

func (h *HealthHandler) GetHealth(w http.ResponseWriter, r *http.Request) {
    log.Printf("🏥 Requisição de health check recebida de: %s", r.RemoteAddr)
    
    health, err := h.healthService.GetHealth()
    if err != nil {
        log.Printf("❌ Erro no health check: %v", err)
        http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
        return
    }
    
    response := domain.Response{
        Message: "API funcionando normalmente",
        Data:    health,
    }
    
    w.Header().Set("Content-Type", "application/json")
    
    if health.Database.Status == "unhealthy" {
        w.WriteHeader(http.StatusServiceUnavailable)
    }
    
    json.NewEncoder(w).Encode(response)
} 