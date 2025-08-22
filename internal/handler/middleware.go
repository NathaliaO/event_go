package handler

import (
    "context"
    "log"
    "net/http"
    "strings"

    "github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(jwtSecret []byte) func(http.HandlerFunc) http.HandlerFunc {
    return func(next http.HandlerFunc) http.HandlerFunc {
        return func(w http.ResponseWriter, r *http.Request) {
            authHeader := r.Header.Get("Authorization")
            if authHeader == "" {
                log.Printf("❌ Tentativa de acesso sem token: %s %s", r.Method, r.URL.Path)
                http.Error(w, "Token não fornecido", http.StatusUnauthorized)
                return
            }

            tokenParts := strings.Split(authHeader, " ")
            if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
                log.Printf("❌ Formato de token inválido: %s %s", r.Method, r.URL.Path)
                http.Error(w, "Formato de token inválido", http.StatusUnauthorized)
                return
            }

            tokenString := tokenParts[1]

            claims := jwt.MapClaims{}
            token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
                return jwtSecret, nil
            })

            if err != nil || !token.Valid {
                log.Printf("❌ Token inválido: %s %s - %v", r.Method, r.URL.Path, err)
                http.Error(w, "Token inválido", http.StatusUnauthorized)
                return
            }

            email := claims["email"].(string)
            log.Printf("✅ Acesso autorizado: %s %s - Usuário: %s", r.Method, r.URL.Path, email)

            ctx := r.Context()
            ctx = context.WithValue(ctx, "user_id", int(claims["user_id"].(float64)))
            ctx = context.WithValue(ctx, "email", email)
            r = r.WithContext(ctx)

            next.ServeHTTP(w, r)
        }
    }
} 