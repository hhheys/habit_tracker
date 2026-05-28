DROP INDEX IF EXISTS outbox_event_status_created_at_idx;

DROP TABLE IF EXISTS outbox_event;

DROP TYPE IF EXISTS event_status;