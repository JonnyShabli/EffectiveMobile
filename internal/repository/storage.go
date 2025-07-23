package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/JonnyShabli/EffectiveMobile/internal/models"
	"github.com/JonnyShabli/EffectiveMobile/pkg/logster"
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type StorageInterface interface {
	InsertSub(ctx context.Context, log logster.Logger, sub *models.Subscription) (string, error)
	GetSub(ctx context.Context, log logster.Logger, subId string) ([]*models.Subscription, error)
	UpdateSub(ctx context.Context, log logster.Logger, sub *models.Subscription) error
	DeleteSub(ctx context.Context, log logster.Logger, sub_id string) error
	ListSub(ctx context.Context, log logster.Logger) ([]*models.Subscription, error)
	SumPriceByDate(ctx context.Context, log logster.Logger, params *models.SumPriceRequest) (int, error)
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

func (s *Storage) GetSub(ctx context.Context, log logster.Logger, subId string) ([]*models.Subscription, error) {
	result := make([]*models.Subscription, 0)
	builder := squirrel.Select("sub_id", "service_name", "price", "user_id", "start_date").
		From("subscriptions").
		Where(squirrel.Eq{"sub_id": subId})
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
	var count int
	sqlstringSelect := `
		SELECT COUNT(*) 
		FROM subscriptions 
		WHERE sub_id = $1 AND deleted_at IS NULL
		`

	builderUpdate := squirrel.Update("subscriptions").
		Set("service_name", sub.Service_name).
		Set("price", sub.Price).
		Set("user_id", sub.User_id).
		Set("start_date", sub.Start_date).
		Set("updated_at", time.Now()).
		Set("deleted_at", sub.Deleted_at).
		Where(squirrel.Eq{"sub_id": sub.Sub_id})

	sqlstringUpdate, argsUpdate, err := builderUpdate.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return err
	}

	log.Infof("begin transaction")
	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}
	log.WithField("sql", sqlstringSelect).WithField("args", sub.Sub_id).Infof("executing query")
	err = tx.QueryRowxContext(ctx, sqlstringSelect, sub.Sub_id).Scan(&count)
	if err != nil {
		errTx := tx.Rollback()
		if errTx != nil {
			return errTx
		}
		return err
	}

	if count != 1 {
		errtx := tx.Rollback()
		if errtx != nil {
			return errtx
		}
	}

	fmt.Println(sqlstringUpdate, argsUpdate)

	log.WithField("sql", sqlstringUpdate).WithField("args", argsUpdate).Infof("executing query")
	_, err = tx.ExecContext(ctx, sqlstringUpdate, argsUpdate...)
	if err != nil {
		err = tx.Rollback()
		if err != nil {
			return err
		}
		return err
	}

	log.Infof("commiting transaction")
	err = tx.Commit()
	if err != nil {
		return err
	}

	return err
}

func (s *Storage) DeleteSub(ctx context.Context, log logster.Logger, subId string) error {
	builder := squirrel.Update("subscriptions").
		Set("deleted_at", time.Now()).
		Where(squirrel.Eq{"sub_id": subId})

	sqlstring, args, err := builder.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return err
	}

	log.WithField("sql", sqlstring).Infof("executing query")

	result, err := s.db.ExecContext(ctx, sqlstring, args...)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return err
}

func (s *Storage) ListSub(ctx context.Context, log logster.Logger) ([]*models.Subscription, error) {
	result := make([]*models.Subscription, 0)
	builder := squirrel.Select("*").
		From("subscriptions").
		Where(squirrel.Eq{"deleted_at": nil})

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

func (s *Storage) SumPriceByDate(ctx context.Context, log logster.Logger, params *models.SumPriceRequest) (int, error) {
	var sum int

	builder := squirrel.Select("COALESCE(SUM(price), 0) as total").
		From("subscriptions").
		Where(squirrel.And{
			squirrel.Expr("start_date BETWEEN ? AND ?", params.Start_date, params.End_date),
			squirrel.Eq{"deleted_at": nil},
		})
	if params.User_id != "" {
		builder = builder.Where(squirrel.Eq{"user_id": params.User_id})
	}
	if params.Service_name != "" {
		builder = builder.Where(squirrel.Eq{"service_name": params.Service_name})
	}

	sqlstring, args, err := builder.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return 0, err
	}

	log.WithField("sql", sqlstring).Infof("executing query")

	err = s.db.GetContext(ctx, &sum, sqlstring, args...)
	if err != nil {
		return 0, err
	}
	return sum, nil
}
