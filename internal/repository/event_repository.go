package repository

import (
	"crypto/sha256"
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/nathaliaoliveira/goapp/internal/domain"
)

type DBInterface interface {
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
}

type eventRepository struct {
	db DBInterface
}

func NewEventRepository(db DBInterface) EventRepository {
	return &eventRepository{db: db}
}

func (r *eventRepository) Create(event *domain.EmailEvent) (string, error) {
	contentHash := r.generateContentHash(event)
	
	var exists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM email_events WHERE content_hash = $1)", contentHash).Scan(&exists)
	if err != nil {
		return "", fmt.Errorf("erro ao verificar duplicação: %w", err)
	}
	
	if exists {
		return "", fmt.Errorf("evento duplicado")
	}
	
	eventID := uuid.New().String()
	
	_, err = r.db.Exec(`
		INSERT INTO email_events (event_id, event_type, email, site, timestamp, content_hash)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, eventID, event.Type, event.Email, event.Site, event.Timestamp, contentHash)
	
	if err != nil {
		return "", fmt.Errorf("erro ao inserir evento: %w", err)
	}
	
	return eventID, nil
}

func (r *eventRepository) generateContentHash(event *domain.EmailEvent) string {
    content := fmt.Sprintf("%s|%s|%s|%s", event.Type, event.Email, event.Site, event.Timestamp)
    hash := sha256.Sum256([]byte(content))
    return fmt.Sprintf("%x", hash)
}

func (r *eventRepository) GetDailyStats(startDate, endDate, site string) ([]domain.DailyStats, error) {
	baseQuery := `
		SELECT 
			DATE(timestamp) as date,
			site,
			event_type,
			COUNT(*) as count,
			COUNT(DISTINCT email) as unique_emails
		FROM email_events
		WHERE 1=1
	`
	
	var args []interface{}
	var conditions []string
	argIndex := 1
	
	if startDate != "" {
		conditions = append(conditions, fmt.Sprintf("DATE(timestamp) >= $%d", argIndex))
		args = append(args, startDate)
		argIndex++
	}
	
	if endDate != "" {
		conditions = append(conditions, fmt.Sprintf("DATE(timestamp) <= $%d", argIndex))
		args = append(args, endDate)
		argIndex++
	}
	
	if site != "" {
		conditions = append(conditions, fmt.Sprintf("site = $%d", argIndex))
		args = append(args, site)
		argIndex++
	}
	
	if len(conditions) > 0 {
		baseQuery += " AND " + strings.Join(conditions, " AND ")
	}
	
	baseQuery += `
		GROUP BY DATE(timestamp), site, event_type
		ORDER BY date DESC, site, event_type
	`
	
	rows, err := r.db.Query(baseQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("erro ao consultar estatísticas: %v", err)
	}
	defer rows.Close()
	
	statsByDate := make(map[string]map[string]map[string]domain.EventStats)
	
	for rows.Next() {
		var date, siteName, eventType string
		var count, uniqueEmails int
		
		err := rows.Scan(&date, &siteName, &eventType, &count, &uniqueEmails)
		if err != nil {
			return nil, fmt.Errorf("erro ao ler dados: %v", err)
		}
        
        if statsByDate[date] == nil {
            statsByDate[date] = make(map[string]map[string]domain.EventStats)
        }
        if statsByDate[date][siteName] == nil {
            statsByDate[date][siteName] = make(map[string]domain.EventStats)
        }
        
        statsByDate[date][siteName][eventType] = domain.EventStats{
            Count:        count,
            UniqueEmails: uniqueEmails,
        }
    }
    
    var result []domain.DailyStats
    
    for date, sites := range statsByDate {
        for siteName, events := range sites {
            totalEvents := 0
            totalUniqueEmails := 0
            
            for _, eventStats := range events {
                totalEvents += eventStats.Count
                totalUniqueEmails += eventStats.UniqueEmails
            }
            
            siteStats := domain.DailyStats{
                Date:              date,
                Site:              siteName,
                TotalEvents:       totalEvents,
                TotalUniqueEmails: totalUniqueEmails,
                Events:            events,
            }
            
            result = append(result, siteStats)
        }
    }
    
    	return result, nil
}

func (r *eventRepository) GetTotalCounts() (int, int, error) {
	var totalUsers, totalEvents int
	
	err := r.db.QueryRow("SELECT COUNT(*) FROM users").Scan(&totalUsers)
	if err != nil {
		return 0, 0, err
	}
	
	err = r.db.QueryRow("SELECT COUNT(*) FROM email_events").Scan(&totalEvents)
	if err != nil {
		return 0, 0, err
	}
	
	return totalUsers, totalEvents, nil
} 