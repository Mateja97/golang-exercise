package storage

import (
	"golang-exercise/models"
)

type storage interface {
	methods
	tx
	BeginTransaction() error
}

type tx interface {
	CommitRollback() error
}

type methods interface {
	CreateCompany(company *models.Company) error
	UpdateCompany(company *models.Company) error
	DeleteCompany(companyID string) error
	ReadCompany(companyID string) (*models.Company, error)
}

var s storage

func Register(i storage) {
	s = i
}

func BeginTransaction() error {
	return s.BeginTransaction()
}

func CreateCompany(company *models.Company) error {
	return s.CreateCompany(company)
}

func UpdateCompany(company *models.Company) error {
	return s.UpdateCompany(company)
}

func ReadCompany(companyID string) (*models.Company, error) {
	return s.ReadCompany(companyID)
}

func DeleteCompany(companyID string) error {
	return s.DeleteCompany(companyID)
}


func CommitRollback() error {
	return s.CommitRollback()
}
