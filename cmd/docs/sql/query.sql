-- show all available metric types
SELECT name FROM app.types;

-- insert data to counter
INSERT INTO
    app.counter_part (metric_id, name, value, created)
VALUES ((SELECT uuid FROM app.types WHERE name = 'counter'), 'test_counter', 22, NOW());

SELECT created FROM  app.counter_part;
SELECT created FROM  app.counter_y2024_3_quarter;

-- insert data to gauge
INSERT INTO
    app.gauge_part (metric_id, name, value)
VALUES ((SELECT uuid FROM app.types WHERE name = 'gauge'), 'test_gauge', 10.010203040506070809);


-- SELECT LAST DATA FROM COUNTER
-- cost=25.73..51.18 rows=4 width=564) (actual time=0.049..0.050 rows=1 loops=1)
EXPLAIN ANALYSE
SELECT
    *
FROM
    app.counter_part
WHERE
    created = (SELECT MAX(created) FROM app.counter_part WHERE name = 'test_counter') AND name = 'test_counter';

-- (cost=0.14..8.16 rows=1 width=540) (actual time=0.056..0.056 rows=0 loops=1)
EXPLAIN ANALYSE
SELECT
    name
FROM
    app.types
WHERE
    uuid = '533af14c-c86f-4ae4-880e-9bb90129c6ef';


-- SELECT LAST DATA for test_counter FROM COUNTER + type
--  (cost=51.35..57.50 rows=1 width=1048) (actual time=0.115..0.116 rows=1 loops=1)
EXPLAIN ANALYSE
WITH LatestCounter AS (
    SELECT
        MAX(created) AS latest_created
    FROM
        app.counter_part
    WHERE
        app.counter_part.name = 'test_counter'
)
SELECT
    t.name, c.name, c.value
FROM
    app.counter_part as c
JOIN
    app.types as t ON c.metric_id = t.uuid
JOIN
    LatestCounter lc ON c.created = lc.latest_created
WHERE
    c.name = 'test_counter'
ORDER BY c.created;

-- Вывод: делать 2 запроса в БД не эффективнее, чтобы собрать все данные, для ответа. А использование 1 запроса, позволит
-- уменьшить кол-во кода в сервисе + 2 запрос в БД вместо одного и не нужно открывать транзакции в коде.

-- SELECT LAST DATA gauge_counter FROM GAUGE + type
--  (cost=51.26..57.40 rows=1 width=1048) (actual time=0.082..0.084 rows=0 loops=1)
EXPLAIN ANALYSE
WITH LatestCounter AS (
    SELECT
        MAX(created) AS latest_created
    FROM
        app.gauge_part
    WHERE
        app.gauge_part.name = 'test_gauge'
)
SELECT
    t.name, g.name, g.value
FROM
    app.gauge_part as g
JOIN
    app.types as t ON g.metric_id = t.uuid
JOIN
    LatestCounter lc ON g.created = lc.latest_created
WHERE
    g.name = 'test_counter'
ORDER BY g.created;

-- SELECT ALL DATA
-- (cost=44.05..83.47 rows=1 width=1040) (actual time=0.100..0.110 rows=3 loops=1)
EXPLAIN ANALYSE
WITH LatestCounter AS (
    SELECT
        MAX(created) AS latest_created, name
    FROM
        app.counter_part
    GROUP BY
        name
)

SELECT
    t.name, c.name, c.value
FROM
    app.counter_part as c
JOIN
    app.types as t ON c.metric_id = t.uuid
JOIN
    LatestCounter lc ON c.created = lc.latest_created AND c.name = lc.name;


-- (cost=43.97..83.31 rows=1 width=1040) (actual time=0.060..0.063 rows=1 loops=1)
EXPLAIN ANALYSE
WITH LatestCounter AS (
    SELECT
        MAX(created) AS latest_created, name
    FROM
        app.gauge_part
    GROUP BY
        name
)

SELECT
    t.name, g.name, g.value
FROM
    app.gauge_part as g
        JOIN
    app.types as t ON g.metric_id = t.uuid
        JOIN
    LatestCounter lc ON g.created = lc.latest_created AND g.name = lc.name



--------------------
-- (cost=166.81..166.83 rows=2 width=1040) (actual time=0.136..0.141 rows=4 loops=1)
EXPLAIN ANALYSE
WITH LatestCounter AS (
    SELECT
        MAX(created) AS latest_created, name
    FROM
        app.counter_part
    GROUP BY
        name
), LatestGauge AS (
    SELECT
        MAX(created) AS latest_created, name
    FROM
        app.gauge_part
    GROUP BY
        name
)
SELECT
    t.name, c.name, c.value
FROM
    app.counter_part as c
        JOIN
    app.types as t ON c.metric_id = t.uuid
        JOIN
    LatestCounter lc ON c.created = lc.latest_created AND c.name = lc.name
UNION
SELECT
    t.name, g.name, g.value
FROM
    app.gauge_part as g
        JOIN
    app.types as t ON g.metric_id = t.uuid
        JOIN
    LatestGauge lc ON g.created = lc.latest_created AND g.name = lc.name