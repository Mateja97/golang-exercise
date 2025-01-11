package api

import (
	"context"
	"fmt"
	companyService "github.com/mateja97/golang-exercise/protobuf/golang/client/service/company/v1"
	"golang-exercise/api/codes"
	"golang-exercise/storage"
	"gorm.io/gorm"
)

func (h *handler) Delete(ctx context.Context, request *companyService.DeleteRequest) (*companyService.DeleteResponse, error) {
	if !IsValidUUID(request.GetId()) {
		fmt.Println("invalid argument")
		return nil, codes.InvalidArgument
	}
	err := storage.Transaction(func(tx *gorm.DB) error {
		if err := storage.DeleteCompany(request.GetId()); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		fmt.Printf("transaction error: %v", err)
		return nil, codes.Internal
	}

	return &companyService.DeleteResponse{}, nil
}
