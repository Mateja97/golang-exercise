package api

import (
	"context"
	companyService "github.com/mateja97/golang-exercise/protobuf/golang/client/service/company/v1"
	"go.uber.org/zap"
	"golang-exercise/api/codes"
	"golang-exercise/logger"
	"golang-exercise/storage"
)

func (h *handler) Get(_ context.Context, request *companyService.GetRequest) (*companyService.GetResponse, error) {
	logger.Info("Get request", zap.String("id", request.GetId()))
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

	return &companyService.GetResponse{
		Company: companyToProto(*company),
	}, nil
}
