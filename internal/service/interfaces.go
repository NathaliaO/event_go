package service

import "github.com/nathaliaoliveira/goapp/internal/domain"

type UserService interface {
    Register(name, email, password string) (*domain.User, error)
    Login(email, password string) (*domain.AuthResponse, error)
    GetByID(id int) (*domain.User, error)
    GetAll() ([]domain.User, error)
    Create(name, email, password string) (*domain.User, error)
}

type EventService interface {
    ProcessEvents(events []domain.EmailEvent) (*domain.EventsResponse, error)
    GetDailyStats(startDate, endDate, site string) (*domain.StatsResponse, error)
}

type HealthService interface {
    GetHealth() (*domain.HealthResponse, error)
} 