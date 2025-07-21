package repository

import (
	"context"
	"time"

	"github.com/JonnyShabli/EffectiveMobile/internal/models"
	"github.com/JonnyShabli/EffectiveMobile/pkg/logster"
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type StorageInterface interface {
	InsertSub(ctx context.Context, log logster.Logger, sub *models.Subscription) (string, error)
	GetSub(ctx context.Context, log logster.Logger, name string, userId string) ([]*models.Subscription, error)
	UpdateSub(ctx context.Context, log logster.Logger, sub *models.Subscription) error
	DeleteSub(ctx context.Context, log logster.Logger, sub *models.Subscription) error
	ListSub(ctx context.Context, log logster.Logger) ([]*models.Subscription, error)
}

type Storage struct {
	db *sqlx.DB
}

func NewStorage(db *sqlx.DB) *Storage {
	return &Storage{db: db}
}

func (s *Storage) InsertSub(ctx context.Context, log logster.Logger, sub *models.Subscription) (string, error) {
	var sub_id string
	builder := squirrel.Insert("subscriptions").
		Columns("service_name", "price", "user_id", "start_date").
		Values(sub.Service_name, sub.Price, sub.User_id, sub.Start_date).
		Suffix("RETURNING \"sub_id\"")

	sqlstring, args, err := builder.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return "", err
	}
	log.WithField("sql", sqlstring).Infof("executing query")
	err = s.db.GetContext(ctx, &sub_id, sqlstring, args...)
	if err != nil {
		return "", err
	}
	return sub_id, nil
}

func (s *Storage) GetSub(ctx context.Context, log logster.Logger, name string, userId string) ([]*models.Subscription, error) {
	result := make([]*models.Subscription, 0)
	builder := squirrel.Select("service_name", "price", "user_id", "start_date").
		From("subscriptions").
		Where(squirrel.Eq{"service_name": name, "user_id": userId})
	sqlstring, args, err := builder.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return nil, err
	}
	log.WithField("sql", sqlstring).Infof("executing query")
	err = s.db.SelectContext(ctx, &result, sqlstring, args...)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *Storage) UpdateSub(ctx context.Context, log logster.Logger, sub *models.Subscription) error {
	builder := squirrel.Update("subscriptions").
		Set("service_name", sub.Service_name).
		Set("price", sub.Price).
		Set("user_id", sub.User_id).
		Set("start_date", sub.Start_date).
		Set("updated_at", time.Now()).
		Where(squirrel.Eq{"service_name": sub.Service_name, "user_id": sub.User_id})

	sqlstring, args, err := builder.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return err
	}

	log.WithField("sql", sqlstring).Infof("executing query")

	_, err = s.db.ExecContext(ctx, sqlstring, args...)
	if err != nil {
		return err
	}

	return err
}

func (s *Storage) DeleteSub(ctx context.Context, log logster.Logger, sub *models.Subscription) error {
	builder := squirrel.Delete("subscriptions").
		Where(squirrel.Eq{"service_name": sub.Service_name, "user_id": sub.User_id})

	sqlstring, args, err := builder.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return err
	}

	log.WithField("sql", sqlstring).Infof("executing query")

	_, err = s.db.ExecContext(ctx, sqlstring, args...)
	if err != nil {
		return err
	}
	return err
}

func (s *Storage) ListSub(ctx context.Context, log logster.Logger) ([]*models.Subscription, error) {
	result := make([]*models.Subscription, 0)
	builder := squirrel.Select("*").
		From("subscriptions")

	sqlstring, args, err := builder.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	log.WithField("sql", sqlstring).Infof("executing query")

	err = s.db.SelectContext(ctx, &result, sqlstring, args...)
	if err != nil {
		return nil, err
	}

	return result, nil
}
