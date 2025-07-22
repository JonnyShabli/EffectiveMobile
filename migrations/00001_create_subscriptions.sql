-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS subscriptions
(
sub_id                  SERIAL PRIMARY KEY,
    service_name        VARCHAR(50)        NOT NULL,
    price               INTEGER            NULL,
    user_id             VARCHAR(50)       NOT NULL,
    start_date          VARCHAR(7)              NOT NULL,

    created_at      TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMP DEFAULT NULL
    );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS subscriptions;
-- +goose StatementEnd