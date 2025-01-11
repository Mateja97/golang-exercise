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

func (h *handler) Patch(_ context.Context, request *companyService.PatchRequest) (*companyService.PatchResponse, error) {
	logger.Info("Patch request", zap.String("id", request.GetId()))
	if !validatePatchRequest(request) {
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

	company, err := storage.ReadCompany(request.GetId())
	if err != nil {
		logger.Error("read company error",
			zap.Error(err))
		return nil, codes.ErrInternal
	}
	if company == nil {
		logger.Info("company not found")
		return nil, codes.ErrNotFound
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
		logger.Error("update company error",
			zap.Error(err))
		return nil, codes.ErrInternal
	}
	event := &models.Event{
		State:     models.StateUpdated,
		CompanyID: company.ID,
		Company:   company,
	}

	err = h.writer.Write(event)
	if err != nil {
		logger.Error("fail to write",
			zap.Error(err))
		return nil, codes.ErrInternal
	}

	return &companyService.PatchResponse{
		Company: companyToProto(*company),
	}, nil
}

func validatePatchRequest(request *companyService.PatchRequest) bool {

	if !isValidUUID(request.Id) {
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

func isValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}
