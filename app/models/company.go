package models

type CompanyType int8

const (
	CompanyTypeUnspecified CompanyType = iota
	CompanyTypeCorporation
	CompanyTypeNonProfit
	CompanyTypeCooperative
	CompanyTypeSoleProprietorship
)

type Company struct {
	ID                string `gorm:"primaryKey"`
	Name              string
	Description       string
	AmountOfEmployees uint32
	Registered        bool
	Type              CompanyType
}
