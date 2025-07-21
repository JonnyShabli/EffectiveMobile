package service

import (
	"context"

	"github.com/JonnyShabli/EffectiveMobile/internal/models"
	"github.com/JonnyShabli/EffectiveMobile/internal/repository"
	"github.com/JonnyShabli/EffectiveMobile/pkg/logster"
)

type SubsServiceInterface interface {
	InsertSub(ctx context.Context, log logster.Logger, sub *models.Subscription) (string, error)
	GetSub(ctx context.Context, log logster.Logger, name string, userId string) ([]*models.Subscription, error)
	UpdateSub(ctx context.Context, log logster.Logger, sub *models.Subscription) error
	DeleteSub(ctx context.Context, log logster.Logger, sub_id string) error
	ListSub(ctx context.Context, log logster.Logger) ([]*models.Subscription, error)
}

type SubsService struct {
	db repository.StorageInterface
}

func NewSubsService(db repository.StorageInterface) *SubsService {
	return &SubsService{db: db}
}

func (s *SubsService) InsertSub(ctx context.Context, log logster.Logger, sub *models.Subscription) (string, error) {
	id, err := s.db.InsertSub(ctx, log, sub)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (s *SubsService) GetSub(ctx context.Context, log logster.Logger, name string, userId string) ([]*models.Subscription, error) {
	subs, err := s.db.GetSub(ctx, log, name, userId)
	if err != nil {
		return nil, err
	}
	return subs, nil
}

func (s *SubsService) UpdateSub(ctx context.Context, log logster.Logger, sub *models.Subscription) error {
	return s.db.UpdateSub(ctx, log, sub)
}

func (s *SubsService) DeleteSub(ctx context.Context, log logster.Logger, sub_id string) error {
	return s.db.DeleteSub(ctx, log, sub_id)
}

func (s *SubsService) ListSub(ctx context.Context, log logster.Logger) ([]*models.Subscription, error) {
	subs, err := s.db.ListSub(ctx, log)
	if err != nil {
		return nil, err
	}
	return subs, nil
}
