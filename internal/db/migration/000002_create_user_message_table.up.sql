CREATE TABLE IF NOT EXISTS user_messages
(
    id         bigserial not null primary key,
    user_id    uuid      not null,
    message    varchar   not null,
    created_at timestamp default current_timestamp,
    updated_at timestamp default null,
    CONSTRAINT fk_user
        FOREIGN KEY (user_id)
            REFERENCES users (id)
            ON DELETE RESTRICT
);