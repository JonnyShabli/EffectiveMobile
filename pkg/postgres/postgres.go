package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/JonnyShabli/EffectiveMobile/pkg/logster"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Config struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Database string `yaml:"database"`

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

func InitDB(db *sqlx.DB) error {
	panic("not implemented")
	return nil
}
