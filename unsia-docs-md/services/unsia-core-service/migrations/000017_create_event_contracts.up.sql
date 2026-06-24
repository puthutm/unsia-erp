-- core_db: 000017_create_event_contracts.up.sql
-- Event contract catalog: registry semua event yang dipublish/dikonsumsi lintas modul

CREATE TABLE event_contracts (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_name              VARCHAR(200) NOT NULL,
    event_version           VARCHAR(10) NOT NULL DEFAULT 'v1',
    event_type              VARCHAR(30) NOT NULL
                            CHECK (event_type IN ('DOMAIN_EVENT', 'INTEGRATION_EVENT', 'NOTIFICATION_EVENT', 'SNAPSHOT_EVENT')),
    publisher_module        VARCHAR(50) NOT NULL,
    publisher_database      VARCHAR(50),
    aggregate_type          VARCHAR(100) NOT NULL,
    payload_schema          JSONB NOT NULL,
    validation_schema       JSONB,
    status                  VARCHAR(20) NOT NULL DEFAULT 'active'
                            CHECK (status IN ('draft', 'active', 'deprecated', 'retired')),
    backward_compatible     BOOLEAN NOT NULL DEFAULT TRUE,
    description             TEXT,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT uq_event_name_version UNIQUE (event_name, event_version)
);

CREATE INDEX idx_event_contracts_publisher_module ON event_contracts (publisher_module);
CREATE INDEX idx_event_contracts_status ON event_contracts (status);
