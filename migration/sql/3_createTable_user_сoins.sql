CREATE TABLE merch.user_coins (
    user_id INTEGER PRIMARY KEY REFERENCES merch.users(id),
    coins BIGINT DEFAULT 0
);

ALTER TABLE merch.user_coins OWNER TO merch_service_user;