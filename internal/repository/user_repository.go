package repository

import (
    "database/sql"
    "log"
    "strings"

    "github.com/nathaliaoliveira/goapp/internal/domain"
)

type userRepository struct {
    db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
    return &userRepository{db: db}
}

func (r *userRepository) Create(name, email, passwordHash string) (*domain.User, error) {
    log.Printf("📝 Criando usuário: %s (%s)", name, email)
    
    query := "INSERT INTO users (name, email, password_hash) VALUES ($1, $2, $3) RETURNING id, name, email, created_at"
    
    var user domain.User
    err := r.db.QueryRow(query, name, email, passwordHash).Scan(
        &user.ID, &user.Name, &user.Email, &user.CreatedAt)
    
    if err != nil {
        if strings.Contains(err.Error(), "unique constraint") {
            log.Printf("❌ Email já cadastrado: %s", email)
            return nil, &DuplicateEmailError{Email: email}
        }
        log.Printf("❌ Erro ao criar usuário: %v", err)
        return nil, err
    }
    
    log.Printf("✅ Usuário criado com sucesso: %s (ID: %d)", user.Name, user.ID)
    return &user, nil
}

func (r *userRepository) GetByID(id int) (*domain.User, error) {
    log.Printf("🔍 Buscando usuário por ID: %d", id)
    
    var user domain.User
    query := "SELECT id, name, email, created_at FROM users WHERE id = $1"
    
    err := r.db.QueryRow(query, id).Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)
    if err != nil {
        if err == sql.ErrNoRows {
            log.Printf("❌ Usuário não encontrado: ID %d", id)
            return nil, &UserNotFoundError{ID: id}
        }
        log.Printf("❌ Erro ao buscar usuário: %v", err)
        return nil, err
    }
    
    log.Printf("✅ Usuário encontrado: %s (ID: %d)", user.Name, user.ID)
    return &user, nil
}

func (r *userRepository) GetByEmail(email string) (*domain.User, error) {
    log.Printf("🔍 Buscando usuário por email: %s", email)
    
    var user domain.User
    query := "SELECT id, name, email, password_hash, created_at FROM users WHERE email = $1"
    
    err := r.db.QueryRow(query, email).Scan(
        &user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.CreatedAt)
    if err != nil {
        if err == sql.ErrNoRows {
            log.Printf("❌ Usuário não encontrado: %s", email)
            return nil, &UserNotFoundError{Email: email}
        }
        log.Printf("❌ Erro ao buscar usuário: %v", err)
        return nil, err
    }
    
    log.Printf("✅ Usuário encontrado: %s (%s)", user.Name, user.Email)
    return &user, nil
}

func (r *userRepository) GetAll() ([]domain.User, error) {
    log.Printf("👥 Listando todos os usuários")
    
    query := "SELECT id, name, email, created_at FROM users ORDER BY id"
    rows, err := r.db.Query(query)
    if err != nil {
        log.Printf("❌ Erro ao listar usuários: %v", err)
        return nil, err
    }
    defer rows.Close()
    
    var users []domain.User
    for rows.Next() {
        var user domain.User
        if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt); err != nil {
            log.Printf("❌ Erro ao ler usuário: %v", err)
            return nil, err
        }
        users = append(users, user)
    }
    
    log.Printf("✅ %d usuários encontrados", len(users))
    return users, nil
}

type DuplicateEmailError struct {
    Email string
}

func (e *DuplicateEmailError) Error() string {
    return "email já cadastrado: " + e.Email
}

type UserNotFoundError struct {
    ID    int
    Email string
}

func (e *UserNotFoundError) Error() string {
    if e.ID != 0 {
        return "usuário não encontrado com ID: " + string(rune(e.ID))
    }
    return "usuário não encontrado com email: " + e.Email
} 