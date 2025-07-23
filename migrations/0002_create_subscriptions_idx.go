package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {

	goose.AddMigrationContext(UP_002, Down_002)
}

func UP_002(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `CREATE INDEX IF NOT EXISTS idx_subscriptions_user_id ON subscriptions (user_id);
CREATE INDEX IF NOT EXISTS idx_subscriptions_service_name ON subscriptions (service_name);`)
	if err != nil {
		return err
	}
	return nil
}

func Down_002(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `DROP INDEX IF EXISTS idx_subscriptions_user_id RESTRICT;
DROP INDEX IF EXISTS idx_subscriptions_service_name RESTRICT`)
	if err != nil {
		return err
	}
	return nil
}
