BEGIN;

CREATE TABLE IF NOT EXISTS urls (
    "id"         bigserial PRIMARY KEY,
    "url"        varchar   NOT NULL,
    "short_uri"  varchar   UNIQUE NOT NULL,
    "views"      int       NOT NULL DEFAULT 0,
    "created_at" timestamptz    NOT NULL DEFAULT (now()),
    "deleted_at" timestamptz NULL
);

COMMIT;