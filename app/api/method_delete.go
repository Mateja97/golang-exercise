package api

import (
	"context"
	companyService "github.com/mateja97/golang-exercise/protobuf/golang/client/service/company/v1"
	"go.uber.org/zap"
	"golang-exercise/api/codes"
	"golang-exercise/logger"
	"golang-exercise/models"
	"golang-exercise/storage"
)

func (h *handler) Delete(_ context.Context, request *companyService.DeleteRequest) (*companyService.DeleteResponse, error) {
	logger.Info("Delete request", zap.String("id", request.GetId()))
	if !isValidUUID(request.GetId()) {
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
	if err := storage.DeleteCompany(request.GetId()); err != nil {
		logger.Error("delete company error",
			zap.Error(err))
		return nil, codes.ErrInternal
	}
	event := &models.Event{
		State:     models.StateDeleted,
		CompanyID: request.GetId(),
		Company:   nil,
	}

	err = h.writer.Write(event)
	if err != nil {
		logger.Error("fail to write",
			zap.Error(err))
		return nil, codes.ErrInternal
	}

	return &companyService.DeleteResponse{}, nil
}
