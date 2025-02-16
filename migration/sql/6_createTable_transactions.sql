CREATE TABLE
    merch.transactions (
        id BIGSERIAL PRIMARY KEY,
        type INTEGER NOT NULL,
        created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        sender_id INTEGER REFERENCES merch.users (id),
        recipient_id INTEGER REFERENCES merch.users (id),
        coins BIGINT DEFAULT 0,
        item_id INTEGER REFERENCES merch.items (id)
    );

ALTER TABLE merch.transactions OWNER TO merch_service_user;

CREATE INDEX idx_transactions_sender_type_recipient ON merch.transactions (sender_id, type, recipient_id) INCLUDE (coins);

CREATE INDEX idx_transactions_recipient_type_sender ON merch.transactions (recipient_id, type, sender_id) INCLUDE (coins);