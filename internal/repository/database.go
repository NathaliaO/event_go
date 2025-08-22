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

func CreateTables(db *sql.DB) error {
    log.Println("📋 Criando tabelas...")
    
    userQuery := `
        CREATE TABLE IF NOT EXISTS users (
            id SERIAL PRIMARY KEY,
            name VARCHAR(100) NOT NULL,
            email VARCHAR(100) UNIQUE NOT NULL,
            password_hash VARCHAR(255) NOT NULL,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        );
    `
    _, err := db.Exec(userQuery)
    if err != nil {
        log.Printf("❌ Erro ao criar tabela users: %v", err)
        return err
    }
    log.Println("✅ Tabela users criada/verificada")

    eventQuery := `
        CREATE TABLE IF NOT EXISTS email_events (
            id SERIAL PRIMARY KEY,
            event_id VARCHAR(100) UNIQUE NOT NULL,
            event_type VARCHAR(50) NOT NULL,
            email VARCHAR(255) NOT NULL,
            site VARCHAR(255) NOT NULL,
            timestamp TIMESTAMP NOT NULL,
            campaign_id VARCHAR(100),
            subject VARCHAR(500),
            ip_address VARCHAR(45),
            user_agent TEXT,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        );
    `
    _, err = db.Exec(eventQuery)
    if err != nil {
        log.Printf("❌ Erro ao criar tabela email_events: %v", err)
        return err
    }
    log.Println("✅ Tabela email_events criada/verificada")

    indexQueries := []string{
        "CREATE INDEX IF NOT EXISTS idx_email_events_email ON email_events(email);",
        "CREATE INDEX IF NOT EXISTS idx_email_events_type ON email_events(event_type);",
        "CREATE INDEX IF NOT EXISTS idx_email_events_timestamp ON email_events(timestamp);",
        "CREATE INDEX IF NOT EXISTS idx_email_events_campaign ON email_events(campaign_id);",
    }

    for i, indexQuery := range indexQueries {
        _, err = db.Exec(indexQuery)
        if err != nil {
            log.Printf("❌ Erro ao criar índice %d: %v", i+1, err)
            return err
        }
    }
    log.Println("✅ Índices criados/verificados")
    
    return nil
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
} 