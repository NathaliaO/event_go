package handler

import (
    "encoding/json"
    "log"
    "net/http"
    "time"

    "github.com/nathaliaoliveira/goapp/internal/domain"
)

type HomeHandler struct{}

func NewHomeHandler() *HomeHandler {
    return &HomeHandler{}
}

func (h *HomeHandler) Home(w http.ResponseWriter, r *http.Request) {
    log.Printf("📄 Página inicial acessada por: %s", r.RemoteAddr)
    
    response := domain.Response{
        Message: "Bem-vindo à API Go com PostgreSQL e JWT!",
        Data: map[string]string{
            "status": "running",
            "time":   time.Now().Format(time.RFC3339),
            "auth_required": "Para acessar rotas protegidas, faça login e use o token JWT",
            "endpoints": "POST /login, POST /register, GET /users (protegido), GET /profile (protegido)",
        },
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
} 