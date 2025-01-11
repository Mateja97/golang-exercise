package api

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	companyService "github.com/mateja97/golang-exercise/protobuf/golang/client/service/company/v1"
	"golang-exercise/api/codes"
	"golang-exercise/models"
	"golang-exercise/storage"
	"gorm.io/gorm"
)

func (h *handler) Get(ctx context.Context, request *companyService.GetRequest) (*companyService.GetResponse, error) {

	if !IsValidUUID(request.GetId()) {
		fmt.Println("invalid argument")
		return nil, codes.InvalidArgument
	}
	var company models.Company
	var err error
	err = storage.Transaction(func(tx *gorm.DB) error {
		if company, err = storage.ReadCompany(request.GetId()); err != nil {
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

	return &companyService.GetResponse{
		Company: companyToProto(company),
	}, nil
}
