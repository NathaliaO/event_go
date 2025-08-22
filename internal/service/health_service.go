package service

import (
	"time"

	"github.com/nathaliaoliveira/goapp/internal/domain"
	"github.com/nathaliaoliveira/goapp/internal/repository"
)

type DBInterface interface {
	Ping() error
}

type healthService struct {
	eventRepo repository.EventRepository
	db        DBInterface
	startTime time.Time
}

func NewHealthService(eventRepo repository.EventRepository, db DBInterface, startTime time.Time) HealthService {
	return &healthService{
		eventRepo: eventRepo,
		db:        db,
		startTime: startTime,
	}
}

func (s *healthService) GetHealth() (*domain.HealthResponse, error) {
	var dbStatus string
	var dbLatency time.Duration
	
	start := time.Now()
	err := s.db.Ping()
	dbLatency = time.Since(start)
	
	if err != nil {
		dbStatus = "unhealthy"
	} else {
		dbStatus = "healthy"
	}
	
	totalUsers, totalEvents, err := s.eventRepo.GetTotalCounts()
	if err != nil {
		totalUsers, totalEvents = 0, 0
	}
	
	uptime := time.Since(s.startTime)
	
	return &domain.HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now().Format(time.RFC3339),
		Database: domain.DatabaseHealth{
			Status:    dbStatus,
			LatencyMs: dbLatency.Milliseconds(),
		},
		Statistics: domain.StatisticsHealth{
			TotalUsers:  totalUsers,
			TotalEvents: totalEvents,
		},
		Uptime: uptime.String(),
	}, nil
} 