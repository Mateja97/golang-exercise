package api

import (
	"context"
	"github.com/google/uuid"
	companyMessage "github.com/mateja97/golang-exercise/protobuf/golang/client/message/company/v1"
	companyService "github.com/mateja97/golang-exercise/protobuf/golang/client/service/company/v1"
	"go.uber.org/zap"
	"golang-exercise/api/codes"
	"golang-exercise/logger"
	"golang-exercise/models"
	"golang-exercise/storage"
)

func (h *handler) Create(_ context.Context, req *companyService.CreateRequest) (*companyService.CreateResponse, error) {
	logger.Info("Create request")

	id := uuid.New().String()

	if !validateCreateRequest(req) {
		logger.Warn("invalid argument")
		return nil, codes.ErrInvalidArgument
	}
	err := storage.BeginTransaction()
	if err != nil {
		logger.Error("fail begin transaction",
			zap.Error(err))
		return nil, codes.ErrInternal
	}

	defer func() {
		err = storage.CommitRollback()
		if err != nil {
			logger.Error("fail commit rollback",
				zap.Error(err))
		}
	}()
	company := companyFromCreateRequest(req, id)
	err = storage.CreateCompany(&company)
	if err != nil {
		logger.Error("create company error",
			zap.Error(err))
		return nil, codes.ErrInternal
	}

	event := &models.Event{
		State:     models.StateCreated,
		CompanyID: company.ID,
		Company:   &company,
	}

	err = h.writer.Write(event)
	if err != nil {
		logger.Error("fail to write",
			zap.Error(err))
		return nil, codes.ErrInternal
	}

	return &companyService.CreateResponse{
		Company: companyToProto(company),
	}, nil
}

func companyFromCreateRequest(req *companyService.CreateRequest, uuid string) models.Company {
	return models.Company{
		ID:                uuid,
		Name:              req.Name,
		Description:       req.Description,
		AmountOfEmployees: req.AmountOfEmployees,
		Registered:        req.Registered,
		Type:              companyTypeFromProto(req.Type),
	}
}
func companyTypeFromProto(companyType companyMessage.CompanyType) models.CompanyType {
	switch companyType {
	case companyMessage.CompanyType_COMPANY_TYPE_CORPORATIONS:
		return models.CompanyTypeCorporation
	case companyMessage.CompanyType_COMPANY_TYPE_NON_PROFIT:
		return models.CompanyTypeNonProfit
	case companyMessage.CompanyType_COMPANY_TYPE_COOPERATIVE:
		return models.CompanyTypeCooperative
	case companyMessage.CompanyType_COMPANY_TYPE_SOLE_PROPRIETORSHIP:
		return models.CompanyTypeSoleProprietorship
	default:
		return models.CompanyTypeUnspecified

	}
}

func validateCreateRequest(req *companyService.CreateRequest) bool {
	if req.Name == "" {
		return false
	}
	if len(req.Name) > 15 {
		return false
	}
	if len(req.Description) > 3000 {
		return false
	}
	if req.Type == companyMessage.CompanyType_COMPANY_TYPE_UNSPECIFIED {
		return false
	}
	return true
}

func companyToProto(company models.Company) *companyMessage.Company {
	return &companyMessage.Company{
		Id:                company.ID,
		Name:              company.Name,
		Description:       company.Description,
		AmountOfEmployees: company.AmountOfEmployees,
		Registered:        company.Registered,
		Type:              companyTypeToProto(company.Type),
	}
}

func companyTypeToProto(companyType models.CompanyType) companyMessage.CompanyType {
	switch companyType {
	case models.CompanyTypeCooperative:
		return companyMessage.CompanyType_COMPANY_TYPE_COOPERATIVE
	case models.CompanyTypeNonProfit:
		return companyMessage.CompanyType_COMPANY_TYPE_NON_PROFIT
	case models.CompanyTypeCorporation:
		return companyMessage.CompanyType_COMPANY_TYPE_CORPORATIONS
	case models.CompanyTypeSoleProprietorship:
		return companyMessage.CompanyType_COMPANY_TYPE_SOLE_PROPRIETORSHIP
	default:
		return companyMessage.CompanyType_COMPANY_TYPE_UNSPECIFIED

	}
}
