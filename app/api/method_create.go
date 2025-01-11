package api

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	companyMessage "github.com/mateja97/golang-exercise/protobuf/golang/client/message/company/v1"
	companyService "github.com/mateja97/golang-exercise/protobuf/golang/client/service/company/v1"
	"golang-exercise/api/codes"
	"golang-exercise/models"
	"golang-exercise/storage"
	"gorm.io/gorm"
)

func (h *handler) Create(ctx context.Context, req *companyService.CreateRequest) (*companyService.CreateResponse, error) {
	id := uuid.New().String()

	if !validateCreateRequest(req) {
		fmt.Println("invalid argument")
		return nil, codes.InvalidArgument
	}
	company := companyFromCreateRequest(req, id)
	err := storage.Transaction(func(tx *gorm.DB) error {
		if err := storage.CreateCompany(company); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		fmt.Printf("transaction error: %v", err)
		return nil, codes.Internal
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
