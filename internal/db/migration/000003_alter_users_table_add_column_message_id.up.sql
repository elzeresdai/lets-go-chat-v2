ALTER TABLE users ADD COLUMN message_id bigint;
ALTER TABLE users
    ADD CONSTRAINT fk_last_message
        FOREIGN KEY (message_id)
            REFERENCES user_messages (id);