package database

import (
	"database/sql"
)

func RunMigrations(db *sql.DB) error {
	
	if err := createTables(db); err != nil {
		return err
	}
	
	if err := createIndexes(db); err != nil {
		return err
	}
	
	return nil
}

func createTables(db *sql.DB) error {
	
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
		return err
	}

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
			content_hash VARCHAR(64),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`
	_, err = db.Exec(eventQuery)
	if err != nil {
		return err
	}
	
	return nil
}

func createIndexes(db *sql.DB) error {
	
	indexQueries := []string{
		"CREATE INDEX IF NOT EXISTS idx_email_events_email ON email_events(email);",
		"CREATE INDEX IF NOT EXISTS idx_email_events_type ON email_events(event_type);",
		"CREATE INDEX IF NOT EXISTS idx_email_events_timestamp ON email_events(timestamp);",
		"CREATE INDEX IF NOT EXISTS idx_email_events_campaign ON email_events(campaign_id);",
	}

	for _, indexQuery := range indexQueries {
		_, err := db.Exec(indexQuery)
		if err != nil {
			return err
		}
	}
	
	return nil
} 