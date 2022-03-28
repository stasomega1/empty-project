package services

import (
	"context"
	"github.com/go-playground/validator/v10"
	"project/inetrnal/app/model"
	"project/inetrnal/app/store"
	"time"

	"github.com/sirupsen/logrus"
)

type ServiceI interface {
	GetSomeData(request model.DbModelRequest) (*model.DbModel, error)
}

type Service struct {
	domain string
	store  *store.Store
	logger *logrus.Logger
}

func NewService(store *store.Store, logger *logrus.Logger, domain string) *Service {
	return &Service{store: store, logger: logger, domain: domain}
}

func getTimeoutContext(timeoutSec int) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Duration(timeoutSec)*time.Second)
}

func (s *Service) GetSomeData(request model.DbModelRequest) (*model.DbModel, error) {
	validate := validator.New()
	err := validate.Struct(request)
	if err != nil {
		return nil, model.NewError("Ошибка валидации запроса", 400, err.Error())
	}

	ctx, cancel := getTimeoutContext(10)
	defer cancel()

	result, err := s.store.SqlRepository().GetSomeDataFromDB(ctx, request.Parameter1, request.Parameter2)
	if err != nil {
		return nil, model.NewError("Серверная ошибка", 500, err.Error())
	}

	return result, nil
}
