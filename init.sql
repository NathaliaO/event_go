CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS email_events (
    id SERIAL PRIMARY KEY,
    event_id VARCHAR(100) UNIQUE NOT NULL,
    event_type VARCHAR(50) NOT NULL,
    email VARCHAR(255) NOT NULL,
    site VARCHAR(255) NOT NULL,
    timestamp TIMESTAMP NOT NULL,
    content_hash VARCHAR(64) UNIQUE NOT NULL,
    campaign_id VARCHAR(100),
    subject VARCHAR(500),
    ip_address VARCHAR(45),
    user_agent TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_email_events_email ON email_events(email);
CREATE INDEX IF NOT EXISTS idx_email_events_type ON email_events(event_type);
CREATE INDEX IF NOT EXISTS idx_email_events_timestamp ON email_events(timestamp);
CREATE INDEX IF NOT EXISTS idx_email_events_campaign ON email_events(campaign_id);
CREATE INDEX IF NOT EXISTS idx_email_events_content_hash ON email_events(content_hash);

-- Eventos de exemplo com content_hash
INSERT INTO email_events (event_id, event_type, email, site, timestamp, content_hash, campaign_id, subject, ip_address, user_agent, created_at) VALUES
('evt_001', 'sent', 'user@example.com', 'site-a.com', '2025-08-21T10:30:00Z', 'sent|user@example.com|site-a.com|2025-08-21T10:30:00Z', 'camp_123', 'Welcome Email', NULL, NULL, '2025-08-21T10:30:00Z'),
('evt_002', 'open', 'user@example.com', 'site-a.com', '2025-08-21T10:35:00Z', 'open|user@example.com|site-a.com|2025-08-21T10:35:00Z', 'camp_123', NULL, '192.168.1.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36', '2025-08-21T10:35:00Z'),
('evt_003', 'click', 'user@example.com', 'site-a.com', '2025-08-21T10:40:00Z', 'click|user@example.com|site-a.com|2025-08-21T10:40:00Z', 'camp_123', NULL, '192.168.1.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36', '2025-08-21T10:40:00Z'),
('evt_004', 'bounce', 'invalid@example.com', 'site-a.com', '2025-08-21T11:00:00Z', 'bounce|invalid@example.com|site-a.com|2025-08-21T11:00:00Z', 'camp_123', 'Welcome Email', NULL, NULL, '2025-08-21T11:00:00Z'),
('evt_005', 'sent', 'user2@example.com', 'site-b.com', '2025-08-21T12:00:00Z', 'sent|user2@example.com|site-b.com|2025-08-21T12:00:00Z', 'camp_456', 'Newsletter Weekly', NULL, NULL, '2025-08-21T12:00:00Z'),
('evt_006', 'open', 'user2@example.com', 'site-b.com', '2025-08-21T12:15:00Z', 'open|user2@example.com|site-b.com|2025-08-21T12:15:00Z', 'camp_456', NULL, '192.168.1.2', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)', '2025-08-21T12:15:00Z'),

('evt_007', 'sent', 'user3@example.com', 'site-c.com', '2025-08-20T09:00:00Z', 'sent|user3@example.com|site-c.com|2025-08-20T09:00:00Z', 'camp_789', 'Promoção Especial', NULL, NULL, '2025-08-20T09:00:00Z'),
('evt_008', 'open', 'user3@example.com', 'site-c.com', '2025-08-20T09:30:00Z', 'open|user3@example.com|site-c.com|2025-08-20T09:30:00Z', 'camp_789', NULL, '192.168.1.3', 'Mozilla/5.0 (iPhone; CPU iPhone OS 14_7_1)', '2025-08-20T09:30:00Z'),
('evt_009', 'click', 'user3@example.com', 'site-c.com', '2025-08-20T09:35:00Z', 'click|user3@example.com|site-c.com|2025-08-20T09:35:00Z', 'camp_789', NULL, '192.168.1.3', 'Mozilla/5.0 (iPhone; CPU iPhone OS 14_7_1)', '2025-08-20T09:35:00Z'),
('evt_010', 'sent', 'user4@example.com', 'site-a.com', '2025-08-20T14:00:00Z', 'sent|user4@example.com|site-a.com|2025-08-20T14:00:00Z', 'camp_123', 'Welcome Email', NULL, NULL, '2025-08-20T14:00:00Z'),

('evt_011', 'sent', 'user5@example.com', 'site-d.com', '2025-08-22T08:00:00Z', 'sent|user5@example.com|site-d.com|2025-08-22T08:00:00Z', 'camp_999', 'Confirmação de Pedido', NULL, NULL, '2025-08-22T08:00:00Z'),
('evt_012', 'open', 'user5@example.com', 'site-d.com', '2025-08-22T08:30:00Z', 'open|user5@example.com|site-d.com|2025-08-22T08:30:00Z', 'camp_999', NULL, '192.168.1.4', 'Mozilla/5.0 (Linux; Android 11)', '2025-08-22T08:30:00Z'),
('evt_013', 'sent', 'user6@example.com', 'site-e.com', '2025-08-22T16:00:00Z', 'sent|user6@example.com|site-e.com|2025-08-22T16:00:00Z', 'camp_777', 'Boas-vindas', NULL, NULL, '2025-08-22T16:00:00Z'),
('evt_014', 'bounce', 'invalid2@example.com', 'site-e.com', '2025-08-22T16:05:00Z', 'bounce|invalid2@example.com|site-e.com|2025-08-22T16:05:00Z', 'camp_777', 'Boas-vindas', NULL, NULL, '2025-08-22T16:05:00Z'),

('evt_015', 'sent', 'user7@example.com', 'site-f.com', '2025-08-19T10:00:00Z', 'sent|user7@example.com|site-f.com|2025-08-19T10:00:00Z', 'camp_555', 'Lembrete de Pagamento', NULL, NULL, '2025-08-19T10:00:00Z'),
('evt_016', 'open', 'user7@example.com', 'site-f.com', '2025-08-19T10:45:00Z', 'open|user7@example.com|site-f.com|2025-08-19T10:45:00Z', 'camp_555', NULL, '192.168.1.5', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64)', '2025-08-19T10:45:00Z'),
('evt_017', 'sent', 'user8@example.com', 'site-g.com', '2025-08-19T15:00:00Z', 'sent|user8@example.com|site-g.com|2025-08-19T15:00:00Z', 'camp_333', 'Confirmação de Cadastro', NULL, NULL, '2025-08-19T15:00:00Z'),

('evt_018', 'sent', 'user9@example.com', 'site-h.com', '2025-08-23T09:00:00Z', 'sent|user9@example.com|site-h.com|2025-08-23T09:00:00Z', 'camp_111', 'Oferta Especial', NULL, NULL, '2025-08-23T09:00:00Z'),
('evt_019', 'open', 'user9@example.com', 'site-h.com', '2025-08-23T09:20:00Z', 'open|user9@example.com|site-h.com|2025-08-23T09:20:00Z', 'camp_111', NULL, '192.168.1.6', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)', '2025-08-23T09:20:00Z'),
('evt_020', 'click', 'user9@example.com', 'site-h.com', '2025-08-23T09:25:00Z', 'click|user9@example.com|site-h.com|2025-08-23T09:25:00Z', 'camp_111', NULL, '192.168.1.6', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)', '2025-08-23T09:25:00Z'); 