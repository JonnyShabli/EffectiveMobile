-- +goose Up
CREATE TABLE IF NOT EXISTS subscriptions
(
    sub_id              BIGINT
    service_name        VARCHAR(50)        NOT NULL,
    price               INTEGER            NULL,
    user_id             VARCHAR(16)        NOT NULL,
    start_date          DATE               NOT NULL,

    created_at      TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMP          DEFAULT NULL,

    PRIMARY KEY (service_name, user_id)
    );
-- +goose Down
DROP TABLE IF EXISTS subscriptions;
