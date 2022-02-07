CREATE TABLE "users"
(
    "id"         uuid PRIMARY KEY,
    "name"       varchar     NOT NULL,
    "hash"       varchar     NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now())
);