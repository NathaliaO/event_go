package handler

import (
    "encoding/json"
    "log"
    "net/http"

    "github.com/nathaliaoliveira/goapp/internal/domain"
    "github.com/nathaliaoliveira/goapp/internal/repository"
    "github.com/nathaliaoliveira/goapp/internal/service"
)

type UserHandler struct {
    userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
    return &UserHandler{
        userService: userService,
    }
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
    log.Printf("📝 Requisição de registro recebida de: %s", r.RemoteAddr)
    
    var registerReq domain.RegisterRequest
    if err := json.NewDecoder(r.Body).Decode(&registerReq); err != nil {
        log.Printf("❌ Dados de registro inválidos: %v", err)
        http.Error(w, "Dados inválidos", http.StatusBadRequest)
        return
    }
    
    user, err := h.userService.Register(registerReq.Name, registerReq.Email, registerReq.Password)
    if err != nil {
        log.Printf("❌ Erro no registro: %v", err)
        h.handleServiceError(w, err)
        return
    }
    
    response := domain.Response{
        Message: "Usuário criado com sucesso",
        Data:    user,
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(response)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
    log.Printf("🔐 Requisição de login recebida de: %s", r.RemoteAddr)
    
    var loginReq domain.LoginRequest
    if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
        log.Printf("❌ Dados de login inválidos: %v", err)
        http.Error(w, "Dados inválidos", http.StatusBadRequest)
        return
    }
    
    authResp, err := h.userService.Login(loginReq.Email, loginReq.Password)
    if err != nil {
        log.Printf("❌ Erro no login: %v", err)
        h.handleServiceError(w, err)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(authResp)
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
    log.Printf("👥 Requisição de listagem de usuários recebida de: %s", r.RemoteAddr)
    
    users, err := h.userService.GetAll()
    if err != nil {
        log.Printf("❌ Erro ao listar usuários: %v", err)
        h.handleServiceError(w, err)
        return
    }
    
    response := domain.Response{
        Message: "Usuários encontrados",
        Data:    users,
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
    log.Printf("👤 Requisição de criação de usuário recebida de: %s", r.RemoteAddr)
    
    var createUserReq struct {
        Name     string `json:"name"`
        Email    string `json:"email"`
        Password string `json:"password"`
    }
    
    if err := json.NewDecoder(r.Body).Decode(&createUserReq); err != nil {
        log.Printf("❌ Dados inválidos para criar usuário: %v", err)
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    user, err := h.userService.Create(createUserReq.Name, createUserReq.Email, createUserReq.Password)
    if err != nil {
        log.Printf("❌ Erro ao criar usuário: %v", err)
        h.handleServiceError(w, err)
        return
    }
    
    response := domain.Response{
        Message: "Usuário criado com sucesso",
        Data:    user,
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(response)
}

func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
    log.Printf("👤 Requisição de perfil recebida de: %s", r.RemoteAddr)
    
    userID := r.Context().Value("user_id").(int)
    
    user, err := h.userService.GetByID(userID)
    if err != nil {
        log.Printf("❌ Erro ao buscar perfil: %v", err)
        h.handleServiceError(w, err)
        return
    }
    
    response := domain.Response{
        Message: "Perfil do usuário",
        Data:    user,
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func (h *UserHandler) handleServiceError(w http.ResponseWriter, err error) {
    switch e := err.(type) {
    case *service.ValidationError:
        http.Error(w, e.Error(), http.StatusBadRequest)
    case *service.AuthenticationError:
        http.Error(w, e.Error(), http.StatusUnauthorized)
    case *service.InternalError:
        http.Error(w, e.Error(), http.StatusInternalServerError)
    default:
        if _, ok := err.(*repository.DuplicateEmailError); ok {
            http.Error(w, err.Error(), http.StatusConflict)
        } else if _, ok := err.(*repository.UserNotFoundError); ok {
            http.Error(w, err.Error(), http.StatusNotFound)
        } else {
            http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
        }
    }
} 