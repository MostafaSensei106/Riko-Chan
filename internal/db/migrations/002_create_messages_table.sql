CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS messages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    recipient_id BIGINT,
    group_id VARCHAR(255),
    channel_id VARCHAR(255),
    message_type VARCHAR(50) NOT NULL,
    content TEXT NOT NULL,
    media_file_id VARCHAR(255),
    location JSONB,
    scheduled_time TIMESTAMP WITH TIME ZONE NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    recurrence_type VARCHAR(20) NOT NULL DEFAULT 'none',
    recurrence_count INTEGER NOT NULL DEFAULT 0,
    max_recurrences INTEGER,
    notify_before INTERVAL,
    private_view_mode BOOLEAN NOT NULL DEFAULT false,
    google_calendar_id VARCHAR(255),
    notion_page_id VARCHAR(255),
    trello_card_id VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);
