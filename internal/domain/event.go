package domain

type EmailEvent struct {
    Type      string                 `json:"type"`
    Email     string                 `json:"email"`
    Site      string                 `json:"site"`
    Timestamp string                 `json:"timestamp"`
    Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

type EventsRequest struct {
    Events []EmailEvent `json:"events"`
}

type ProcessedEvent struct {
    ID        string `json:"id"`
    Type      string `json:"type"`
    Email     string `json:"email"`
    Site      string `json:"site"`
    Status    string `json:"status"` // "processed", "duplicate", "error"
}

type EventsResponse struct {
    Processed  int              `json:"processed"`
    Duplicates int              `json:"duplicates"`
    Errors     int              `json:"errors"`
    Events     []ProcessedEvent `json:"events"`
} 