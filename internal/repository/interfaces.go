package repository

import "github.com/nathaliaoliveira/goapp/internal/domain"

type UserRepository interface {
    Create(name, email, passwordHash string) (*domain.User, error)
    GetByEmail(email string) (*domain.User, error)
    GetByID(id int) (*domain.User, error)
    GetAll() ([]domain.User, error)
}

type EventRepository interface {
    Create(event *domain.EmailEvent) (string, error)
    GetDailyStats(startDate, endDate, site string) ([]domain.DailyStats, error)
    GetTotalCounts() (int, int, error)
} 