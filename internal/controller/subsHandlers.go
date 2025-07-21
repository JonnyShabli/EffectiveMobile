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
		BadRequestResponse(w, "fail read request body", nil, err.Error())
		return
	}
	err = json.Unmarshal(reqBody, &subDTO)
	if err != nil {
		BadRequestResponse(w, "json.Unmarshal", subDTO, err.Error())
		return
	}
	sub := dtoToSub(&subDTO)

	sub_id, err := s.Service.InsertSub(ctx, s.Log, sub)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			err = fmt.Errorf("duplicate key value %w", err)
			BadRequestResponse(w, "duplicate key value", nil, err.Error())
			return
		}

		ErrorResponse(w, "Service.InsertSub error", nil, err.Error())
		return
	}

	SuccessResponse(w, "success", sub_id, "")
}

func (s *SubsHandler) GetSub(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	serviceName := chi.URLParam(r, "service_name")
	userId := chi.URLParam(r, "user_id")
	if serviceName == "" && userId == "" {
		BadRequestResponse(w, "bad params", nil, "")
		return
	}
	subs, err := s.Service.GetSub(ctx, s.Log, serviceName, userId)
	if err != nil {
		if strings.Contains(err.Error(), "sql: no rows in result set") {
			err = fmt.Errorf("sql: no rows in result set %w", err)
			BadRequestResponse(w, "sql: no rows in result set", nil, err.Error())
			return
		}

		ErrorResponse(w, "Service.GetSub error", nil, err.Error())
		return
	}
	result := subsToDTO(subs)

	SuccessResponse(w, "success", result, "")
}

func (s *SubsHandler) UpdateSub(w http.ResponseWriter, r *http.Request) {
	var subDTO models.SubscriptionDTO
	ctx := r.Context()
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		BadRequestResponse(w, "fail read request body", nil, err.Error())
		return
	}
	err = json.Unmarshal(reqBody, &subDTO)
	if err != nil {
		BadRequestResponse(w, "json.Unmarshal", subDTO, err.Error())
		return
	}
	sub := dtoToSub(&subDTO)

	err = s.Service.UpdateSub(ctx, s.Log, sub)
	if err != nil {
		BadRequestResponse(w, "Service.UpdateSub error", sub, err.Error())
		return
	}
	SuccessResponse(w, "success", nil, "")
}

func (s *SubsHandler) DeleteSub(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sub_id := chi.URLParam(r, "sub_id")
	if sub_id == "" {
		BadRequestResponse(w, "bad params", nil, "")
		return
	}
	err := s.Service.DeleteSub(ctx, s.Log, sub_id)
	if err != nil {
		BadRequestResponse(w, "Service.DeleteSub error", sub_id, err.Error())
		return
	}
	SuccessResponse(w, "success", nil, "")
}

func (s *SubsHandler) ListSub(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	subs, err := s.Service.ListSub(ctx, s.Log)
	if err != nil {
		if strings.Contains(err.Error(), "sql: no rows in result set") {
			err = fmt.Errorf("sql: no rows in result set %w", err)
			BadRequestResponse(w, "sql: no rows in result set", nil, err.Error())
			return
		}

		ErrorResponse(w, "Service.GetSub error", nil, err.Error())
		return
	}
	result := subsToDTO(subs)

	SuccessResponse(w, "success", result, "")

}

func dtoToSub(dto *models.SubscriptionDTO) *models.Subscription {
	return &models.Subscription{
		Price:        dto.Price,
		Service_name: dto.Service_name,
		User_id:      dto.User_id,
		Start_date:   dto.Start_date,
	}
}

func subsToDTO(subs []*models.Subscription) []*models.SubscriptionDTO {
	result := make([]*models.SubscriptionDTO, len(subs))
	for _, sub := range subs {
		result = append(result, &models.SubscriptionDTO{
			Price:        sub.Price,
			Service_name: sub.Service_name,
			User_id:      sub.User_id,
			Start_date:   sub.Start_date,
		})
	}
	return result
}
