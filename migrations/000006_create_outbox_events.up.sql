CREATE TYPE event_status AS ENUM ('created', 'sent', 'dead');

CREATE TABLE IF NOT EXISTS outbox_event (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        occurred_at TIMESTAMP NOT NULL,
        event_type VARCHAR(256) NOT NULL,
        event_type_version INTEGER NOT NULL DEFAULT 1,
        status event_status DEFAULT 'created' NOT NULL,
        partition_key VARCHAR(128) NOT NULL,
        attempt_count INTEGER NOT NULL DEFAULT 0,
        next_attempt_at TIMESTAMP NOT NULL DEFAULT NOW(),
        payload JSONB NOT NULL,
        created_at TIMESTAMP NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS outbox_event_status_created_at_idx
    ON outbox_event (status, created_at);

CREATE INDEX IF NOT EXISTS outbox_event_status_next_attempt_idx
    ON outbox_event (status, next_attempt_at, created_at);