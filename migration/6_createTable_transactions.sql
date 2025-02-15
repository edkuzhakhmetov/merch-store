CREATE TABLE merch.transactions (
    id BIGSERIAL PRIMARY KEY,
    type INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    sender_id INTEGER REFERENCES merch.users(id),
    recipient_id INTEGER REFERENCES merch.users(id),
    coins BIGINT DEFAULT 0, 
    item_id INTEGER REFERENCES merch.items(id)
);

ALTER TABLE merch.transactions OWNER TO merch_service_user;