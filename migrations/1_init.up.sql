CREATE SCHEMA IF NOT EXISTS app;

CREATE TABLE IF NOT EXISTS app.types (
    uuid uuid PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    created TIMESTAMP WITH TIME ZONE,
    modified TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS app.metrics (
    uuid uuid PRIMARY KEY,
    metric_id uuid NOT NULL REFERENCES app.types(uuid) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL UNIQUE,
    value DOUBLE PRECISION,
    created TIMESTAMP WITH TIME ZONE,
    modified TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS metrics_modified_idx
    ON app.metrics (modified);
