package domain

type Response struct {
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}

type Claims struct {
    UserID int    `json:"user_id"`
    Email  string `json:"email"`
    Exp    int64  `json:"exp"`
    Iat    int64  `json:"iat"`
}

type HealthResponse struct {
    Status    string                 `json:"status"`
    Timestamp string                 `json:"timestamp"`
    Database  DatabaseHealth         `json:"database"`
    Statistics StatisticsHealth      `json:"statistics"`
    Uptime    string                 `json:"uptime"`
}

type DatabaseHealth struct {
    Status     string `json:"status"`
    LatencyMs  int64  `json:"latency_ms"`
}

type StatisticsHealth struct {
    TotalUsers  int `json:"total_users"`
    TotalEvents int `json:"total_events"`
} 