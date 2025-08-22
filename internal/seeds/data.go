package seeds

import (
	"database/sql"

	"golang.org/x/crypto/bcrypt"
)

func RunSeeds(db *sql.DB) error {
	if err := createInitialUsers(db); err != nil {
		return err
	}
	
	if err := createSampleEvents(db); err != nil {
		return err
	}
	
	return nil
}

func createInitialUsers(db *sql.DB) error {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		return err
	}
	
	if count > 0 {
		return nil
	}
	
	adminPassword := "admin123"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	
	_, err = db.Exec(`
		INSERT INTO users (name, email, password_hash) 
		VALUES ($1, $2, $3)
		ON CONFLICT (email) DO NOTHING
	`, "Admin User", "admin@test.com", string(hashedPassword))
	
	if err != nil {
		return err
	}
	
	return nil
}

func createSampleEvents(db *sql.DB) error {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM email_events").Scan(&count)
	if err != nil {
		return err
	}
	
	if count > 0 {
		return nil
	}
	
	sampleEvents := []struct {
		eventID, eventType, email, site, timestamp, contentHash, campaignID, subject, ipAddress, userAgent string
	}{
		{"evt_001", "sent", "user@example.com", "site-a.com", "2025-08-21T10:30:00Z", "sent|user@example.com|site-a.com|2025-08-21T10:30:00Z", "camp_123", "Welcome Email", "", ""},
		{"evt_002", "open", "user@example.com", "site-a.com", "2025-08-21T10:35:00Z", "open|user@example.com|site-a.com|2025-08-21T10:35:00Z", "camp_123", "", "192.168.1.1", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"},
		{"evt_003", "click", "user@example.com", "site-a.com", "2025-08-21T10:40:00Z", "click|user@example.com|site-a.com|2025-08-21T10:40:00Z", "camp_123", "", "192.168.1.1", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"},
		{"evt_004", "bounce", "invalid@example.com", "site-a.com", "2025-08-21T11:00:00Z", "bounce|invalid@example.com|site-a.com|2025-08-21T11:00:00Z", "camp_123", "Welcome Email", "", ""},
		{"evt_005", "sent", "user2@example.com", "site-b.com", "2025-08-21T12:00:00Z", "sent|user2@example.com|site-b.com|2025-08-21T12:00:00Z", "camp_456", "Newsletter Weekly", "", ""},
		{"evt_006", "open", "user2@example.com", "site-b.com", "2025-08-21T12:15:00Z", "open|user2@example.com|site-b.com|2025-08-21T12:15:00Z", "camp_456", "", "192.168.1.2", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)"},
		{"evt_007", "sent", "user3@example.com", "site-c.com", "2025-08-20T09:00:00Z", "sent|user3@example.com|site-c.com|2025-08-20T09:00:00Z", "camp_789", "Promoção Especial", "", ""},
		{"evt_008", "open", "user3@example.com", "site-c.com", "2025-08-20T09:30:00Z", "open|user3@example.com|site-c.com|2025-08-20T09:30:00Z", "camp_789", "", "192.168.1.3", "Mozilla/5.0 (iPhone; CPU iPhone OS 14_7_1)"},
		{"evt_009", "click", "user3@example.com", "site-c.com", "2025-08-20T09:35:00Z", "click|user3@example.com|site-c.com|2025-08-20T09:35:00Z", "camp_789", "", "192.168.1.3", "Mozilla/5.0 (iPhone; CPU iPhone OS 14_7_1)"},
		{"evt_010", "sent", "user4@example.com", "site-a.com", "2025-08-20T14:00:00Z", "sent|user4@example.com|site-a.com|2025-08-20T14:00:00Z", "camp_123", "Welcome Email", "", ""},
	}
	
	for _, event := range sampleEvents {
		_, err = db.Exec(`
			INSERT INTO email_events (event_id, event_type, email, site, timestamp, content_hash, campaign_id, subject, ip_address, user_agent)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
			ON CONFLICT (event_id) DO NOTHING
		`, event.eventID, event.eventType, event.email, event.site, event.timestamp, event.contentHash, event.campaignID, event.subject, event.ipAddress, event.userAgent)
		
		if err != nil {
			return err
		}
	}
	
	return nil
} 