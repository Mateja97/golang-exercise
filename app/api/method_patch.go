package api

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	companyMessage "github.com/mateja97/golang-exercise/protobuf/golang/client/message/company/v1"
	companyService "github.com/mateja97/golang-exercise/protobuf/golang/client/service/company/v1"
	"golang-exercise/api/codes"
	"golang-exercise/models"
	"golang-exercise/storage"
	"gorm.io/gorm"
)

func (h *handler) Patch(ctx context.Context, request *companyService.PatchRequest) (*companyService.PatchResponse, error) {

	if !validatePatchRequest(request) {
		fmt.Println("invalid argument")
		return nil, codes.InvalidArgument
	}
	var company models.Company
	var err error
	err = storage.Transaction(func(tx *gorm.DB) error {
		if company, err = storage.ReadCompany(request.GetId()); err != nil {
			return err
		}
		if request.Description != nil {
			company.Description = request.GetDescription()
		}
		if request.AmountOfEmployees != nil {
			company.AmountOfEmployees = request.GetAmountOfEmployees()
		}
		if request.Registered != nil {
			company.Registered = request.GetRegistered()
		}
		if request.Type != nil {
			company.Type = companyTypeFromProto(request.GetType())
		}

		if err = storage.UpdateCompany(company); err != nil {
			return err
		}
		return nil
	})
	if errors.Is(err, sql.ErrNoRows) {
		fmt.Println("company not found")
		return nil, codes.NotFound
	}
	if err != nil {
		fmt.Printf("transaction error: %v", err)
		return nil, codes.Internal
	}

	return &companyService.PatchResponse{
		Company: companyToProto(company),
	}, nil
}

func validatePatchRequest(request *companyService.PatchRequest) bool {

	if !IsValidUUID(request.Id) {
		return false
	}
	if request.Description != nil && len(request.GetDescription()) > 3000 {
		return false
	}
	if request.Type != nil && *request.Type == companyMessage.CompanyType_COMPANY_TYPE_UNSPECIFIED {
		return false
	}
	return true
}

func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}
