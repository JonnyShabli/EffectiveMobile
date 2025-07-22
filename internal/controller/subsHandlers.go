package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/JonnyShabli/EffectiveMobile/internal/models"
	"github.com/JonnyShabli/EffectiveMobile/internal/service"
	"github.com/JonnyShabli/EffectiveMobile/pkg/logster"
	"github.com/go-chi/chi/v5"
)

type SubsHandlerInterface interface {
	InsertSub(w http.ResponseWriter, r *http.Request)
	GetSub(w http.ResponseWriter, r *http.Request)
	UpdateSub(w http.ResponseWriter, r *http.Request)
	DeleteSub(w http.ResponseWriter, r *http.Request)
	ListSub(w http.ResponseWriter, r *http.Request)
}

func NewSubsHandler(service service.SubsServiceInterface, log logster.Logger) SubsHandlerInterface {
	return &SubsHandler{
		Service: service,
		Log:     log,
	}
}

type SubsHandler struct {
	Service service.SubsServiceInterface
	Log     logster.Logger
}

func (s *SubsHandler) InsertSub(w http.ResponseWriter, r *http.Request) {
	var subDTO models.SubscriptionDTO
	ctx := r.Context()
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		s.Log.WithError(err).Infof("fail read request body")
		BadRequestResponse(w, s.Log, "fail read request body", "")
		return
	}
	err = json.Unmarshal(reqBody, &subDTO)
	if err != nil {
		err = fmt.Errorf("fail to unmarshal request body '%w'", err)
		s.Log.WithError(err).Infof("fail to unmarshal request body")
		BadRequestResponse(w, s.Log, err.Error(), "")
		return
	}
	sub := dtoToSub(&subDTO)

	sub_id, err := s.Service.InsertSub(ctx, s.Log, sub)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			err = fmt.Errorf("duplicate key value %w", err)
			BadRequestResponse(w, s.Log, err.Error(), "")
			s.Log.WithError(err).Infof("duplicate key value")
			return
		}
		s.Log.WithError(err).Errorf("insert sub fail")
		ErrorResponse(w, s.Log, fmt.Errorf("insert sub fail %w", err).Error(), "")
		return
	}

	SuccessResponse(w, s.Log, "success", sub_id)
	s.Log.Infof("insert sub success")
}

func (s *SubsHandler) GetSub(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	subId := chi.URLParam(r, "sub_id")
	if subId == "" {
		s.Log.Infof("fail to get sub_id")
		BadRequestResponse(w, s.Log, fmt.Errorf("fail to get sub_id").Error(), "")
		return
	}
	subs, err := s.Service.GetSub(ctx, s.Log, subId)
	if err != nil {
		if strings.Contains(err.Error(), "sql: no rows in result set") {
			err = fmt.Errorf("sql: no rows in result set %w", err)
			s.Log.WithError(err).Infof("sql: no rows in result set")
			BadRequestResponse(w, s.Log, err.Error(), "")
			return
		}
		s.Log.WithError(err).Errorf("get sub fail")
		ErrorResponse(w, s.Log, fmt.Errorf("get sub fail %w", err).Error(), "")
		return
	}
	result := subsToDTO(subs)

	SuccessResponse(w, s.Log, "success", result)
	s.Log.Infof("get sub success")
}

func (s *SubsHandler) UpdateSub(w http.ResponseWriter, r *http.Request) {
	var subDTO models.SubscriptionDTO
	ctx := r.Context()
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		err = fmt.Errorf("fail to read request body '%w'", err)
		s.Log.WithError(err).Infof("fail to read request body")
		BadRequestResponse(w, s.Log, err.Error(), "")
		return
	}
	err = json.Unmarshal(reqBody, &subDTO)
	if err != nil {
		err = fmt.Errorf("fail to unmarshal request body '%w'", err)
		s.Log.WithError(err).Infof("fail to unmarshal request body")
		BadRequestResponse(w, s.Log, err.Error(), "")
		return
	}
	sub := dtoToSub(&subDTO)

	err = s.Service.UpdateSub(ctx, s.Log, sub)
	if err != nil {
		err = fmt.Errorf("update sub fail %w", err)
		s.Log.WithError(err).Infof("update sub fail")
		BadRequestResponse(w, s.Log, err.Error(), "")
		return
	}
	SuccessResponse(w, s.Log, "success", "")
	s.Log.Infof("update sub success")
}

func (s *SubsHandler) DeleteSub(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	subId := chi.URLParam(r, "sub_id")
	if subId == "" {
		err := fmt.Errorf("fail to get sub_id")
		s.Log.WithError(err).Infof("fail to get sub_id")
		BadRequestResponse(w, s.Log, err.Error(), "")
		return
	}
	err := s.Service.DeleteSub(ctx, s.Log, subId)
	if err != nil {
		err = fmt.Errorf("delete sub fail %w", err)
		s.Log.WithError(err).Infof("delete sub fail")
		BadRequestResponse(w, s.Log, err.Error(), "")
		return
	}
	SuccessResponse(w, s.Log, "success", "")
	s.Log.Infof("delete sub success")
}

func (s *SubsHandler) ListSub(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	subs, err := s.Service.ListSub(ctx, s.Log)
	if err != nil {
		if strings.Contains(err.Error(), "sql: no rows in result set") {
			err = fmt.Errorf("sql: no rows in result set %w", err)
			BadRequestResponse(w, s.Log, fmt.Errorf("sql: no rows in result set %w", err).Error(), "")
			return
		}
		s.Log.WithError(err).Errorf("list sub fail")
		ErrorResponse(w, s.Log, fmt.Errorf("list sub fail %w", err).Error(), "")
		return
	}
	result := subsToDTO(subs)

	SuccessResponse(w, s.Log, "success", result)
	s.Log.Infof("list sub success")

}

func dtoToSub(dto *models.SubscriptionDTO) *models.Subscription {
	return &models.Subscription{
		Sub_id:       dto.Sub_id,
		Price:        dto.Price,
		Service_name: dto.Service_name,
		User_id:      dto.User_id,
		Start_date:   dto.Start_date,
	}
}

func subsToDTO(subs []*models.Subscription) []*models.SubscriptionDTO {
	result := make([]*models.SubscriptionDTO, 0, len(subs))
	for _, sub := range subs {
		result = append(result, &models.SubscriptionDTO{
			Sub_id:       sub.Sub_id,
			Price:        sub.Price,
			Service_name: sub.Service_name,
			User_id:      sub.User_id,
			Start_date:   sub.Start_date,
		})
	}

	return result
}
