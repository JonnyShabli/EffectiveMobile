package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/JonnyShabli/EffectiveMobile/migrations"
	"github.com/JonnyShabli/EffectiveMobile/pkg/logster"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

type Config struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Database string `yaml:"database"`
	Dialect  string `yaml:"dialect"`
	Migrate  bool   `yaml:"migrate"`

	MinConnections    int           `yaml:"minConnections"`
	MaxConnections    int           `yaml:"maxConnections"`
	HealthCheckPeriod time.Duration `yaml:"healthcheckPeriod"`
	MaxConnIdleTime   time.Duration `yaml:"maxConnIdletime"`
	MaxConnLifetime   time.Duration `yaml:"maxConnLifetime"`
}

func NewConn(ctx context.Context, log logster.Logger, config Config) *sqlx.DB {
	connString := fmt.Sprintf("dbname=%s user=%s password=%s host=%s port=%s sslmode=disable",
		config.Database,
		config.User,
		config.Password,
		config.Host,
		config.Port)

	db, err := sqlx.ConnectContext(ctx, "postgres", connString)
	if err != nil {
		log.WithError(err).Fatalf("failed to connect to postgres: %v", connString)
		return nil
	}

	db.SetMaxIdleConns(config.MaxConnections)
	db.SetMaxOpenConns(config.MaxConnections)
	db.SetConnMaxLifetime(config.MaxConnLifetime)
	db.SetConnMaxIdleTime(config.MaxConnIdleTime)

	return db
}

func MigrateDB(ctx context.Context, db *sql.DB, logger logster.Logger) error {
	goose.SetLogger(logger)
	logger.Infof("apply migrations")
	err := goose.UpContext(ctx, db, ".")

	if err != nil {
		return err
	}

	return nil
}
