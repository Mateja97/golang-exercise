package storage

import (
	"golang-exercise/models"
	"gorm.io/gorm"
)

type storage interface {
	methods
	Transaction(func(tx *gorm.DB) error) error
}

type methods interface {
	CreateCompany(company models.Company) error
	UpdateCompany(company models.Company) error
	DeleteCompany(companyID string) error
	ReadCompany(companyID string) (models.Company, error)
}

var s storage

func Register(i storage) {
	s = i
}

func Transaction(f func(tx *gorm.DB) error) error {
	return s.Transaction(f)
}

func CreateCompany(company models.Company) error {
	return s.CreateCompany(company)
}

func UpdateCompany(company models.Company) error {
	return s.UpdateCompany(company)
}

func ReadCompany(companyID string) (models.Company, error) {
	return s.ReadCompany(companyID)
}

func DeleteCompany(companyID string) error {
	return s.DeleteCompany(companyID)
}
