package pg_storage

import (
	_ "github.com/lib/pq"
	"golang-exercise/models"
	"golang-exercise/storage"
	"gorm.io/gorm"
)

type postgresStorage struct {
	db *gorm.DB
}

var s *postgresStorage

func init() {
	s = &postgresStorage{}
}

func Init(db *gorm.DB) error {
	s.db = db
	storage.Register(s)

	err := db.AutoMigrate(&models.Company{})
	if err != nil {
		return err
	}
	return nil
}

func (p postgresStorage) CreateCompany(company models.Company) error {
	return p.db.Create(company).Error
}

func (p postgresStorage) UpdateCompany(company models.Company) error {
	updates := map[string]interface{}{
		"Description":       company.Description,
		"AmountOfEmployees": company.AmountOfEmployees,
		"Registered":        company.Registered,
		"Type":              company.Type,
	}
	return p.db.Model(&company).Updates(updates).Error
}

func (p postgresStorage) DeleteCompany(companyID string) error {
	return p.db.Delete(&models.Company{}, "id = ?", companyID).Error
}

func (p postgresStorage) ReadCompany(companyID string) (models.Company, error) {
	var company models.Company
	err := p.db.First(&company, "id = ?", companyID).Error
	return company, err
}

func (p postgresStorage) Transaction(f func(tx *gorm.DB) error) error {
	return p.db.Transaction(f)
}
