package service

import (
	"testing"
	"time"

	"github.com/nathaliaoliveira/goapp/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockEventRepository é um mock do repositório para testes
type MockEventRepositoryForHealth struct {
	mock.Mock
}

func (m *MockEventRepositoryForHealth) Create(event *domain.EmailEvent) (string, error) {
	args := m.Called(event)
	return args.String(0), args.Error(1)
}

func (m *MockEventRepositoryForHealth) GetDailyStats(startDate, endDate, site string) ([]domain.DailyStats, error) {
	args := m.Called(startDate, endDate, site)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.DailyStats), args.Error(1)
}

func (m *MockEventRepositoryForHealth) GetTotalCounts() (int, int, error) {
	args := m.Called()
	return args.Int(0), args.Int(1), args.Error(2)
}

// MockDB é um mock simples do banco de dados
type MockDB struct {
	pingError error
}

func (m *MockDB) Ping() error {
	return m.pingError
}

func TestGetHealth_RepositoryError(t *testing.T) {
	mockRepo := new(MockEventRepositoryForHealth)
	mockDB := &MockDB{pingError: nil} // DB saudável
	startTime := time.Now().Add(-time.Hour)
	service := NewHealthService(mockRepo, mockDB, startTime)

	// Mock: repositório falha
	mockRepo.On("GetTotalCounts").Return(0, 0, assert.AnError)

	result, err := service.GetHealth()

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "healthy", result.Status)
	assert.Equal(t, "healthy", result.Database.Status)
	// Mesmo com erro no repositório, deve retornar 0s
	assert.Equal(t, 0, result.Statistics.TotalUsers)
	assert.Equal(t, 0, result.Statistics.TotalEvents)

	mockRepo.AssertExpectations(t)
}

func TestGetHealth_UptimeCalculation(t *testing.T) {
	mockRepo := new(MockEventRepositoryForHealth)
	mockDB := &MockDB{pingError: nil}
	startTime := time.Now().Add(-2 * time.Hour) // 2 horas atrás
	service := NewHealthService(mockRepo, mockDB, startTime)

	// Mock: repositório retorna estatísticas
	mockRepo.On("GetTotalCounts").Return(1, 10, nil)

	result, err := service.GetHealth()

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Contains(t, result.Uptime, "2h")

	mockRepo.AssertExpectations(t)
} 