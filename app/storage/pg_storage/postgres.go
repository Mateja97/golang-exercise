package pgstorage

import (
	"errors"
	"github.com/jinzhu/gorm"
	"golang-exercise/models"
	"golang-exercise/storage"
)

type postgresStorage struct {
	db *gorm.DB
	tx *gorm.DB
}

var ErrNoActiveTransaction = errors.New("no active transaction to commit or rollback")

var s *postgresStorage

func Init(db *gorm.DB) error {
	s = new(postgresStorage)

	s.db = db
	storage.Register(s)

	err := db.AutoMigrate(new(models.Company)).Error
	if err != nil {
		return err
	}
	return nil
}

func (p *postgresStorage) CreateCompany(company *models.Company) error {
	return p.tx.Create(company).Error
}

func (p *postgresStorage) UpdateCompany(company *models.Company) error {
	updates := map[string]interface{}{
		"Description":       company.Description,
		"AmountOfEmployees": company.AmountOfEmployees,
		"Registered":        company.Registered,
		"Type":              company.Type,
	}
	return p.tx.Model(company).Updates(updates).Error

}

func (p *postgresStorage) DeleteCompany(companyID string) error {
	return p.tx.Delete(&models.Company{}, "id = ?", companyID).Error
}

func (p *postgresStorage) ReadCompany(companyID string) (*models.Company, error) {
	company := new(models.Company)

	err := p.tx.First(company, "id = ?", companyID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return company, nil
}

func (p *postgresStorage) BeginTransaction() error {
	p.tx = p.db.Begin()
	return p.tx.Error
}

func (p *postgresStorage) CommitRollback() error {
	if p.tx == nil {
		return ErrNoActiveTransaction
	}

	if r := recover(); r != nil {
		if err := p.tx.Rollback().Error; err != nil {
			return err
		}
	} else if p.db.Error == nil {
		if err := p.tx.Commit().Error; err != nil {
			return err
		}
	} else {
		if err := p.tx.Rollback().Error; err != nil {
			return err
		}
	}
	return nil
}
