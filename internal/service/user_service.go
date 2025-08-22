package service

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nathaliaoliveira/goapp/internal/domain"
	"github.com/nathaliaoliveira/goapp/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
    userRepo   repository.UserRepository
    jwtSecret  []byte
}

func NewUserService(userRepo repository.UserRepository, jwtSecret []byte) UserService {
    return &userService{
        userRepo:  userRepo,
        jwtSecret: jwtSecret,
    }
}

func (s *userService) Register(name, email, password string) (*domain.User, error) {
	if name == "" || email == "" || password == "" {
		return nil, &ValidationError{Message: "Nome, email e senha são obrigatórios"}
	}
	
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, &InternalError{Message: "Erro ao processar senha", Cause: err}
	}
	
	user, err := s.userRepo.Create(name, email, string(hashedPassword))
	if err != nil {
		return nil, err
	}
	
	return user, nil
}

func (s *userService) Login(email, password string) (*domain.AuthResponse, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return nil, &AuthenticationError{Message: "Credenciais inválidas"}
	}
	
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, &AuthenticationError{Message: "Credenciais inválidas"}
	}
	
	token, err := s.generateJWT(user)
	if err != nil {
		return nil, &InternalError{Message: "Erro ao gerar token", Cause: err}
	}
	
	return &domain.AuthResponse{
		Token: token,
		User:  *user,
	}, nil
}

func (s *userService) GetByID(id int) (*domain.User, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	
	return user, nil
}

func (s *userService) GetAll() ([]domain.User, error) {
	users, err := s.userRepo.GetAll()
	if err != nil {
		return nil, err
	}
	
	return users, nil
}

func (s *userService) Create(name, email, password string) (*domain.User, error) {
	if name == "" || email == "" || password == "" {
		return nil, &ValidationError{Message: "Nome, email e senha são obrigatórios"}
	}
	
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, &InternalError{Message: "Erro ao criptografar senha", Cause: err}
	}
	
	user, err := s.userRepo.Create(name, email, string(hashedPassword))
	if err != nil {
		return nil, err
	}
	
	return user, nil
}

func (s *userService) generateJWT(user *domain.User) (string, error) {
    claims := jwt.MapClaims{
        "user_id": user.ID,
        "email":   user.Email,
        "exp":     time.Now().Add(24 * time.Hour).Unix(), // 24 horas
        "iat":     time.Now().Unix(),
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(s.jwtSecret)
}

type ValidationError struct {
    Message string
}

func (e *ValidationError) Error() string {
    return e.Message
}

type AuthenticationError struct {
    Message string
}

func (e *AuthenticationError) Error() string {
    return e.Message
}

type InternalError struct {
    Message string
    Cause   error
}

func (e *InternalError) Error() string {
    if e.Cause != nil {
        return e.Message + ": " + e.Cause.Error()
    }
    return e.Message
} 