-- show all available metric types
SELECT name FROM app.types;

-- insert data to counter
INSERT INTO
    app.counter_part (metric_id, name, value, created)
VALUES ((SELECT uuid FROM app.types WHERE name = 'counter'), 'test_counter', 19, NOW());

SELECT created FROM  app.counter_part;
SELECT created FROM  app.counter_y2024_3_quarter;

-- insert data to gauge
INSERT INTO
    app.gauge_part (metric_id, name, value)
VALUES ((SELECT uuid FROM app.types WHERE name = 'gauge'), 'test_gauge', 10.010203040506070809);

SELECT * FROM  app.gauge_part;
SELECT * FROM  app.gauge_y2024_3_quarter;

-- SELECT LAST DATA FROM COUNTER
SELECT * FROM  app.counter_part WHERE created = (SELECT MAX(created) FROM app.counter_part) AND name = 'test_counter';
SELECT * FROM  app.gauge_part WHERE created = (SELECT MAX(created) FROM app.gauge_part) AND name = 'test_gauge';

