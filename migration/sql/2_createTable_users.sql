CREATE TABLE merch.users (
        id SERIAL PRIMARY KEY,
        username VARCHAR(150) UNIQUE NOT NULL,
        hashed_password VARCHAR(128) NOT NULL,
        salt VARCHAR(64),
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

ALTER TABLE merch.users OWNER TO merch_service_user;

CREATE INDEX IX_UsersName ON merch.users USING hash (username text_pattern_ops)
WITH (fillfactor=90);