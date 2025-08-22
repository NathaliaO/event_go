package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func (c *DatabaseConfig) Connect() (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
		
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar ao banco: %v", err)
	}
	
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("erro ao fazer ping no banco: %v", err)
	}
	
	return db, nil
} 