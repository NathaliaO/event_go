package main

import (
    "crypto/rand"
    "log"
    "net/http"
    "os"
    "time"

    "github.com/gorilla/mux"
    "github.com/nathaliaoliveira/goapp/internal/handler"
    "github.com/nathaliaoliveira/goapp/internal/repository"
    "github.com/nathaliaoliveira/goapp/internal/service"
)

func main() {
    startTime := time.Now()
    
    jwtSecret := generateJWTSecret()
    log.Printf("JWT Secret gerado: %x", jwtSecret)

    dbConfig := repository.NewDatabaseConfig()
    db, err := dbConfig.Connect()
    if err != nil {
        log.Fatal("❌ Erro ao conectar ao banco:", err)
    }
    defer db.Close()

    if err := repository.CreateTables(db); err != nil {
        log.Fatal("❌ Erro ao criar tabelas:", err)
    }

    userRepo := repository.NewUserRepository(db)
    eventRepo := repository.NewEventRepository(db)

    userService := service.NewUserService(userRepo, jwtSecret)
    eventService := service.NewEventService(eventRepo)
    healthService := service.NewHealthService(eventRepo, db, startTime)

    homeHandler := handler.NewHomeHandler()
    userHandler := handler.NewUserHandler(userService)
    eventHandler := handler.NewEventHandler(eventService)
    healthHandler := handler.NewHealthHandler(healthService)

    r := mux.NewRouter()
    
    r.HandleFunc("/", homeHandler.Home).Methods("GET")
    r.HandleFunc("/health", healthHandler.GetHealth).Methods("GET")
    r.HandleFunc("/login", userHandler.Login).Methods("POST")
    r.HandleFunc("/register", userHandler.Register).Methods("POST")
    
    r.HandleFunc("/users", handler.AuthMiddleware(jwtSecret)(userHandler.GetUsers)).Methods("GET")
    r.HandleFunc("/users", handler.AuthMiddleware(jwtSecret)(userHandler.CreateUser)).Methods("POST")
    r.HandleFunc("/profile", handler.AuthMiddleware(jwtSecret)(userHandler.GetProfile)).Methods("GET")
    
    r.HandleFunc("/api/events", handler.AuthMiddleware(jwtSecret)(eventHandler.CreateEvents)).Methods("POST")
    
    r.HandleFunc("/api/stats/daily", handler.AuthMiddleware(jwtSecret)(eventHandler.GetDailyStats)).Methods("GET")

    port := getEnv("PORT", "8080")
    log.Printf("🚀 Servidor rodando na porta %s", port)
    log.Fatal(http.ListenAndServe(":"+port, r))
}

func generateJWTSecret() []byte {
    secret := make([]byte, 32)
    _, err := rand.Read(secret)
    if err != nil {
        log.Fatal("❌ Erro ao gerar JWT secret:", err)
    }
    return secret
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
} 