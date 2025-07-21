-- +goose Up
CREATE INDEX IF NOT EXISTS idx_subscriptions_user_id ON subscriptions (user_id);
CREATE INDEX IF NOT EXISTS idx_subscriptions_service_name ON subscriptions (service_name);


-- +goose Down
DROP INDEX IF EXISTS idx_subscriptions_user_id RESTRICT;
DROP INDEX IF EXISTS idx_subscriptions_service_name RESTRICT

