package service

import (
	"strings"

	"github.com/nathaliaoliveira/goapp/internal/domain"
	"github.com/nathaliaoliveira/goapp/internal/repository"
)

type eventService struct {
    eventRepo repository.EventRepository
}

func NewEventService(eventRepo repository.EventRepository) EventService {
    return &eventService{
        eventRepo: eventRepo,
    }
}

func (s *eventService) ProcessEvents(events []domain.EmailEvent) (*domain.EventsResponse, error) {
	if len(events) == 0 {
		return nil, &ValidationError{Message: "Lista de eventos não pode estar vazia"}
	}
	
	processedCount := 0
	duplicatesCount := 0
	errorsCount := 0
	var processedEvents []domain.ProcessedEvent
	
	for _, event := range events {
		if event.Type == "" || event.Email == "" || event.Site == "" || event.Timestamp == "" {
			errorsCount++
			processedEvents = append(processedEvents, domain.ProcessedEvent{
				Type:   event.Type,
				Email:  event.Email,
				Site:   event.Site,
				Status: "error",
			})
			continue
		}
		
		eventID, err := s.eventRepo.Create(&event)
		if err != nil {
			// Verificar se é erro de duplicação
			if strings.Contains(err.Error(), "evento duplicado") {
				duplicatesCount++
				processedEvents = append(processedEvents, domain.ProcessedEvent{
					Type:   event.Type,
					Email:  event.Email,
					Site:   event.Site,
					Status: "duplicate",
				})
				continue
			}
			
			errorsCount++
			processedEvents = append(processedEvents, domain.ProcessedEvent{
				Type:   event.Type,
				Email:  event.Email,
				Site:   event.Site,
				Status: "error",
			})
			continue
		}
		
		processedCount++
		processedEvents = append(processedEvents, domain.ProcessedEvent{
			ID:     eventID,
			Type:   event.Type,
			Email:  event.Email,
			Site:   event.Site,
			Status: "processed",
		})
	}
	
	return &domain.EventsResponse{
		Processed:  processedCount,
		Duplicates: duplicatesCount,
		Errors:     errorsCount,
		Events:     processedEvents,
	}, nil
}

func (s *eventService) GetDailyStats(startDate, endDate, site string) (*domain.StatsResponse, error) {
	stats, err := s.eventRepo.GetDailyStats(startDate, endDate, site)
	if err != nil {
		return nil, err
	}
	
	return &domain.StatsResponse{
		Period: map[string]string{
			"start_date": startDate,
			"end_date":   endDate,
		},
		SiteFilter: site,
		TotalDays:  len(stats),
		TotalSites: len(stats),
		Stats:      stats,
	}, nil
} 