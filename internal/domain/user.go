package domain

import "time"

type User struct {
    ID           int       `json:"id"`
    Name         string    `json:"name"`
    Email        string    `json:"email"`
    PasswordHash string    `json:"-"`
    CreatedAt    time.Time `json:"created_at"`
}

type LoginRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

type RegisterRequest struct {
    Name     string `json:"name"`
    Email    string `json:"email"`
    Password string `json:"password"`
}

type AuthResponse struct {
    Token string `json:"token"`
    User  User   `json:"user"`
} 