CREATE TABLE merch.items (
    id SERIAL PRIMARY KEY,
    name VARCHAR(250) NOT NULL,
    price bigint NOT NULL DEFAULT 0
);

ALTER TABLE IF EXISTS merch.items
    OWNER to merch_service_user;

INSERT INTO merch.items (name, price) VALUES
    ('t-shirt', 80),
    ('cup', 20),
    ('book', 50),
    ('pen', 10),
    ('powerbank', 200),
    ('hoody', 300),
    ('umbrella', 200),
    ('socks', 10),
    ('wallet', 50),
    ('pink-hoody', 500);
