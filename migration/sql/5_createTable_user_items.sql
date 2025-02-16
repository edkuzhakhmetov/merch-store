CREATE TABLE merch.user_items (
    user_id INTEGER REFERENCES merch.users(id), 
    item_id INTEGER REFERENCES merch.items(id), 
    quantity INTEGER DEFAULT 0, 
	PRIMARY KEY (user_id, item_id)
);

ALTER TABLE merch.user_items OWNER TO merch_service_user;

