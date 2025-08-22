package repository

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func NewDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "postgres"),
		DBName:   getEnv("DB_NAME", "goapp"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}
}

func (c *DatabaseConfig) Connect() (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
	
	log.Printf("🔌 Conectando ao banco: %s:%s/%s", c.Host, c.Port, c.DBName)
	
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar ao banco: %v", err)
	}
	
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("erro ao fazer ping no banco: %v", err)
	}
	
	log.Println("✅ Conectado ao banco de dados PostgreSQL!")
	return db, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
} 