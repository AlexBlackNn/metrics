CREATE SCHEMA IF NOT EXISTS app;


CREATE TABLE IF NOT EXISTS app.types (
                                         uuid uuid PRIMARY KEY DEFAULT gen_random_uuid(),
                                         name VARCHAR(255) NOT NULL UNIQUE,
                                         created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS app.counter_part (
                                           uuid uuid NOT NULL DEFAULT gen_random_uuid(),
                                           metric_id uuid NOT NULL REFERENCES app.types(uuid) ON DELETE CASCADE,
                                           name VARCHAR(255) NOT NULL,
                                           value BIGINT,
                                           created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                           PRIMARY KEY (uuid, created)
) PARTITION BY RANGE (created);

CREATE TABLE IF NOT EXISTS app.counter_y2024_1_quarter  PARTITION OF app.counter_part
    FOR VALUES FROM ('2024-01-01 00:00:00') TO ('2024-04-01 00:00:00');

CREATE TABLE IF NOT EXISTS app.counter_y2024_2_quarter  PARTITION OF app.counter_part
    FOR VALUES FROM ('2024-04-01 00:00:00') TO ('2024-07-01 00:00:00');

CREATE TABLE IF NOT EXISTS app.counter_y2024_3_quarter  PARTITION OF app.counter_part
    FOR VALUES FROM ('2024-07-01 00:00:00') TO ('2024-10-01 00:00:00');

CREATE TABLE IF NOT EXISTS app.counter_y2024_4_quarter  PARTITION OF app.counter_part
    FOR VALUES FROM ('2024-10-01 00:00:00') TO ('2025-01-01 00:00:00');


CREATE TABLE IF NOT EXISTS app.gauge_part (
                                              uuid uuid NOT NULL DEFAULT gen_random_uuid(),
                                              metric_id uuid NOT NULL REFERENCES app.types(uuid) ON DELETE CASCADE,
                                              name VARCHAR(255) NOT NULL,
                                              value DOUBLE PRECISION,
                                              created TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                              PRIMARY KEY (uuid, created)
) PARTITION BY RANGE (created);

CREATE TABLE IF NOT EXISTS app.gauge_y2024_1_quarter  PARTITION OF app.gauge_part
    FOR VALUES FROM ('2024-01-01 00:00:00') TO ('2024-04-01 00:00:00');

CREATE TABLE IF NOT EXISTS app.gauge_y2024_2_quarter  PARTITION OF app.gauge_part
    FOR VALUES FROM ('2024-04-01 00:00:00') TO ('2024-07-01 00:00:00');

CREATE TABLE IF NOT EXISTS app.gauge_y2024_3_quarter  PARTITION OF app.gauge_part
    FOR VALUES FROM ('2024-07-01 00:00:00') TO ('2024-10-01 00:00:00');

CREATE TABLE IF NOT EXISTS app.gauge_y2024_4_quarter  PARTITION OF app.gauge_part
    FOR VALUES FROM ('2024-10-01 00:00:00') TO ('2025-01-01 00:00:00');

-- -- Для каждой дочерней таблицы создаем индекс по ключевому столбцу
CREATE INDEX IF NOT EXISTS counter_created_1_quarter_idx ON app.counter_y2024_1_quarter (created);
CREATE INDEX IF NOT EXISTS counter_created_2_quarter_idx ON app.counter_y2024_2_quarter (created);
CREATE INDEX IF NOT EXISTS counter_created_3_quarter_idx ON app.counter_y2024_3_quarter (created);
CREATE INDEX IF NOT EXISTS counter_created_4_quarter_idx ON app.counter_y2024_4_quarter (created);

CREATE INDEX IF NOT EXISTS counter_name_1_quarter_idx ON app.counter_y2024_1_quarter (name);
CREATE INDEX IF NOT EXISTS counter_name_2_quarter_idx ON app.counter_y2024_2_quarter (name);
CREATE INDEX IF NOT EXISTS counter_name_3_quarter_idx ON app.counter_y2024_3_quarter (name);
CREATE INDEX IF NOT EXISTS counter_name_4_quarter_idx ON app.counter_y2024_4_quarter (name);


-- Для каждой дочерней таблицы создаем индекс по ключевому столбцу
CREATE INDEX IF NOT EXISTS gauge_created_1_quarter_idx ON app.gauge_y2024_1_quarter (created);
CREATE INDEX IF NOT EXISTS gauge_created_2_quarter_idx ON app.gauge_y2024_2_quarter (created);
CREATE INDEX IF NOT EXISTS gauge_created_3_quarter_idx ON app.gauge_y2024_3_quarter (created);
CREATE INDEX IF NOT EXISTS gauge_created_4_quarter_idx ON app.gauge_y2024_4_quarter (created);

CREATE INDEX IF NOT EXISTS gauge_name_1_quarter_idx ON app.gauge_y2024_1_quarter (name);
CREATE INDEX IF NOT EXISTS gauge_name_2_quarter_idx ON app.gauge_y2024_2_quarter (name);
CREATE INDEX IF NOT EXISTS gauge_name_3_quarter_idx ON app.gauge_y2024_3_quarter (name);
CREATE INDEX IF NOT EXISTS gauge_name_4_quarter_idx ON app.gauge_y2024_4_quarter (name);

INSERT INTO app.types(name) VALUES ('gauge');
INSERT INTO app.types(name) VALUES ('counter');