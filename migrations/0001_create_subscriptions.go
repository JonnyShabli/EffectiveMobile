package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {

	goose.AddMigrationContext(UP_001, Down_001)
}

func UP_001(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS subscriptions
(
sub_id                  SERIAL PRIMARY KEY,
    service_name        VARCHAR(50)        NOT NULL,
    price               INTEGER            NULL,
    user_id             VARCHAR(50)       NOT NULL,
    start_date          VARCHAR(7)              NOT NULL,

    created_at      TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMP DEFAULT NULL
    );`)
	if err != nil {
		return err
	}
	return nil
}

func Down_001(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `DROP TABLE IF EXISTS subscriptions;`)
	if err != nil {
		return err
	}
	return nil
}
