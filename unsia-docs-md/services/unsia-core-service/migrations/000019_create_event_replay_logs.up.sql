-- core_db: 000019_create_event_replay_logs.up.sql

CREATE TABLE event_replay_logs (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_key       VARCHAR(500) NOT NULL,
    event_name      VARCHAR(200) NOT NULL,
    event_version   VARCHAR(10),
    source_module   VARCHAR(50) NOT NULL,
    consumer_module VARCHAR(50),
    replay_reason   TEXT NOT NULL,
    replayed_by     UUID,
    replayed_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    replay_status   VARCHAR(20) NOT NULL
                    CHECK (replay_status IN ('queued', 'success', 'failed')),
    last_error      TEXT,
    audit_ref_id    UUID
);

CREATE INDEX idx_event_replay_logs_event_key ON event_replay_logs (event_key);
CREATE INDEX idx_event_replay_logs_event_name ON event_replay_logs (event_name);
CREATE INDEX idx_event_replay_logs_source_module ON event_replay_logs (source_module);
CREATE INDEX idx_event_replay_logs_consumer_module ON event_replay_logs (consumer_module);
CREATE INDEX idx_event_replay_logs_replay_status ON event_replay_logs (replay_status);
CREATE INDEX idx_event_replay_logs_replayed_at ON event_replay_logs (replayed_at);
