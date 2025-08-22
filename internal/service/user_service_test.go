package service

import (
	"testing"

	"github.com/nathaliaoliveira/goapp/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(name, email, passwordHash string) (*domain.User, error) {
	args := m.Called(name, email, passwordHash)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(email string) (*domain.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) GetByID(id int) (*domain.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) GetAll() ([]domain.User, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.User), args.Error(1)
}

func TestRegister_ValidUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo, []byte("test-secret"))

	expectedUser := &domain.User{
		ID:    1,
		Name:  "Test User",
		Email: "test@example.com",
	}

	mockRepo.On("Create", "Test User", "test@example.com", mock.AnythingOfType("string")).Return(expectedUser, nil)

	result, err := service.Register("Test User", "test@example.com", "password123")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedUser.Name, result.Name)
	assert.Equal(t, expectedUser.Email, result.Email)

	mockRepo.AssertExpectations(t)
}

func TestRegister_EmptyFields(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo, []byte("test-secret"))

	result, err := service.Register("", "test@example.com", "password123")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "Nome, email e senha são obrigatórios")

	mockRepo.AssertNotCalled(t, "Create")
}

func TestRegister_RepositoryError(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo, []byte("test-secret"))

	mockRepo.On("Create", "Test User", "test@example.com", mock.AnythingOfType("string")).Return(nil, assert.AnError)

	result, err := service.Register("Test User", "test@example.com", "password123")

	assert.Error(t, err)
	assert.Nil(t, result)

	mockRepo.AssertExpectations(t)
}

func TestLogin_ValidCredentials(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo, []byte("test-secret"))

	hashedPassword := "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi"
	user := &domain.User{
		ID:           1,
		Name:         "Test User",
		Email:        "test@example.com",
		PasswordHash: hashedPassword,
	}

	mockRepo.On("GetByEmail", "test@example.com").Return(user, nil)

	result, err := service.Login("test@example.com", "password")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result.Token)
	assert.Equal(t, user.Email, result.User.Email)

	mockRepo.AssertExpectations(t)
}

func TestLogin_InvalidEmail(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo, []byte("test-secret"))

	mockRepo.On("GetByEmail", "invalid@example.com").Return(nil, assert.AnError)

	result, err := service.Login("invalid@example.com", "password")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "Credenciais inválidas")

	mockRepo.AssertExpectations(t)
}

func TestLogin_InvalidPassword(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo, []byte("test-secret"))

	hashedPassword := "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi"
	user := &domain.User{
		ID:           1,
		Name:         "Test User",
		Email:        "test@example.com",
		PasswordHash: hashedPassword,
	}

	mockRepo.On("GetByEmail", "test@example.com").Return(user, nil)

	result, err := service.Login("test@example.com", "wrongpassword")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "Credenciais inválidas")

	mockRepo.AssertExpectations(t)
}

func TestGetByID_ValidID(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo, []byte("test-secret"))

	expectedUser := &domain.User{
		ID:    1,
		Name:  "Test User",
		Email: "test@example.com",
	}

	mockRepo.On("GetByID", 1).Return(expectedUser, nil)

	result, err := service.GetByID(1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedUser.ID, result.ID)
	assert.Equal(t, expectedUser.Name, result.Name)

	mockRepo.AssertExpectations(t)
}

func TestGetAll_ValidUsers(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo, []byte("test-secret"))

	expectedUsers := []domain.User{
		{ID: 1, Name: "User 1", Email: "user1@example.com"},
		{ID: 2, Name: "User 2", Email: "user2@example.com"},
	}

	mockRepo.On("GetAll").Return(expectedUsers, nil)

	result, err := service.GetAll()

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.Equal(t, expectedUsers[0].Name, result[0].Name)
	assert.Equal(t, expectedUsers[1].Name, result[1].Name)

	mockRepo.AssertExpectations(t)
} 