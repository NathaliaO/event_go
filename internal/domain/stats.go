package domain

type DailyStats struct {
    Date                string                 `json:"date"`
    Site                string                 `json:"site"`
    TotalEvents         int                    `json:"total_events"`
    TotalUniqueEmails   int                    `json:"total_unique_emails"`
    Events              map[string]EventStats  `json:"events"`
}

type EventStats struct {
    Count        int `json:"count"`
    UniqueEmails int `json:"unique_emails"`
}

type StatsRequest struct {
    StartDate string `json:"start_date"`
    EndDate   string `json:"end_date"`
    Site      string `json:"site"`
}

type StatsResponse struct {
    Period     map[string]string `json:"period"`
    SiteFilter string            `json:"site_filter"`
    TotalDays  int               `json:"total_days"`
    TotalSites int               `json:"total_sites"`
    Stats      []DailyStats      `json:"stats"`
} 