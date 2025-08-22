package service

import (
	"testing"

	"github.com/nathaliaoliveira/goapp/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockEventRepository struct {
	mock.Mock
}

func (m *MockEventRepository) Create(event *domain.EmailEvent) (string, error) {
	args := m.Called(event)
	return args.String(0), args.Error(1)
}

func (m *MockEventRepository) GetDailyStats(startDate, endDate, site string) ([]domain.DailyStats, error) {
	args := m.Called(startDate, endDate, site)
	return args.Get(0).([]domain.DailyStats), args.Error(1)
}

func (m *MockEventRepository) GetTotalCounts() (int, int, error) {
	args := m.Called()
	return args.Int(0), args.Int(1), args.Error(2)
}

func TestProcessEvents_ValidEvents(t *testing.T) {
	mockRepo := new(MockEventRepository)
	service := NewEventService(mockRepo)

	events := []domain.EmailEvent{
		{
			Type:      "sent",
			Email:     "user@example.com",
			Site:      "site-a.com",
			Timestamp: "2025-08-20T10:30:00Z",
		},
		{
			Type:      "open",
			Email:     "user@example.com",
			Site:      "site-a.com",
			Timestamp: "2025-08-20T10:35:00Z",
		},
	}

	mockRepo.On("Create", &events[0]).Return("uuid-1", nil)
	mockRepo.On("Create", &events[1]).Return("uuid-2", nil)

	result, err := service.ProcessEvents(events)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 2, result.Processed)
	assert.Equal(t, 0, result.Duplicates)
	assert.Equal(t, 0, result.Errors)
	assert.Len(t, result.Events, 2)
	assert.Equal(t, "processed", result.Events[0].Status)
	assert.Equal(t, "processed", result.Events[1].Status)

	mockRepo.AssertExpectations(t)
}

func TestProcessEvents_DuplicateEvent(t *testing.T) {
	mockRepo := new(MockEventRepository)
	service := NewEventService(mockRepo)

	events := []domain.EmailEvent{
		{
			Type:      "sent",
			Email:     "user@example.com",
			Site:      "site-a.com",
			Timestamp: "2025-08-20T10:30:00Z",
		},
	}

	mockRepo.On("Create", &events[0]).Return("", assert.AnError)

	result, err := service.ProcessEvents(events)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 0, result.Processed)
	assert.Equal(t, 0, result.Duplicates)
	assert.Equal(t, 1, result.Errors)
	assert.Len(t, result.Events, 1)
	assert.Equal(t, "error", result.Events[0].Status)

	mockRepo.AssertExpectations(t)
}

func TestProcessEvents_InvalidEvent(t *testing.T) {
	// Arrange
	mockRepo := new(MockEventRepository)
	service := NewEventService(mockRepo)

	events := []domain.EmailEvent{
		{
			Type:      "", // Campo obrigatório vazio
			Email:     "user@example.com",
			Site:      "site-a.com",
			Timestamp: "2025-08-20T10:30:00Z",
		},
	}

	// Act
	result, err := service.ProcessEvents(events)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 0, result.Processed)
	assert.Equal(t, 0, result.Duplicates)
	assert.Equal(t, 1, result.Errors)
	assert.Len(t, result.Events, 1)
	assert.Equal(t, "error", result.Events[0].Status)

	// Não deve chamar o repositório para eventos inválidos
	mockRepo.AssertNotCalled(t, "Create")
}

func TestProcessEvents_EmptyEventsList(t *testing.T) {
	// Arrange
	mockRepo := new(MockEventRepository)
	service := NewEventService(mockRepo)

	events := []domain.EmailEvent{}

	// Act
	result, err := service.ProcessEvents(events)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "Lista de eventos não pode estar vazia")

	mockRepo.AssertNotCalled(t, "Create")
}

func TestProcessEvents_MixedValidAndInvalid(t *testing.T) {
	// Arrange
	mockRepo := new(MockEventRepository)
	service := NewEventService(mockRepo)

	events := []domain.EmailEvent{
		{
			Type:      "sent",
			Email:     "user@example.com",
			Site:      "site-a.com",
			Timestamp: "2025-08-20T10:30:00Z",
		},
		{
			Type:      "", // Inválido
			Email:     "user@example.com",
			Site:      "site-a.com",
			Timestamp: "2025-08-20T10:35:00Z",
		},
		{
			Type:      "open",
			Email:     "user@example.com",
			Site:      "site-a.com",
			Timestamp: "2025-08-20T10:40:00Z",
		},
	}

	// Mock: primeiro evento processado com sucesso
	mockRepo.On("Create", &events[0]).Return("uuid-1", nil)
	// Mock: terceiro evento processado com sucesso
	mockRepo.On("Create", &events[2]).Return("uuid-3", nil)

	// Act
	result, err := service.ProcessEvents(events)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 2, result.Processed)
	assert.Equal(t, 0, result.Duplicates)
	assert.Equal(t, 1, result.Errors)
	assert.Len(t, result.Events, 3)

	// Verificar status de cada evento
	assert.Equal(t, "processed", result.Events[0].Status)
	assert.Equal(t, "error", result.Events[1].Status)
	assert.Equal(t, "processed", result.Events[2].Status)

	mockRepo.AssertExpectations(t)
} 