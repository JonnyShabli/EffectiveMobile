-- +goose Up
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS idx_subscriptions_user_id on subscriptions (user_id);
CREATE INDEX IF NOT EXISTS idx_subscriptions_service_name on subscriptions (service_name);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF NOT EXISTS idx_subscriptions_user_id RESTRICT;
DROP INDEX IF NOT EXISTS idx_subscriptions_service_name RESTRICT
-- +goose StatementEnd
