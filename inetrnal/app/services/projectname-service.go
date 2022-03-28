package services

import (
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"project/inetrnal/app/model"
	"project/inetrnal/app/store"
	"project/inetrnal/app/utils"
)

type ProjectnameServiceI interface {
	GetSomeData(request model.DbModelRequest) (*model.DbModel, error)
}

type ProjectnameService struct {
	domain string
	store  store.StoreI
	logger *logrus.Logger
}

func NewProjectnameServiceService(store *store.Store, logger *logrus.Logger, domain string) *ProjectnameService {
	return &ProjectnameService{store: store, logger: logger, domain: domain}
}

func (s *ProjectnameService) GetSomeData(request model.DbModelRequest) (*model.DbModel, error) {
	validate := validator.New()
	err := validate.Struct(request)
	if err != nil {
		return nil, model.NewError("Ошибка валидации запроса", 400, err.Error())
	}

	ctx, cancel := utils.GetTimeoutContext(10)
	defer cancel()

	result, err := s.store.SqlRepository().GetSomeDataFromDB(ctx, request.Parameter1, request.Parameter2)
	if err != nil {
		return nil, model.NewError("Серверная ошибка", 500, err.Error())
	}

	return result, nil
}
