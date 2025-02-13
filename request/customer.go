package request

import "time"

type PageCustomer struct {
	Paging
	CompanyID   string     `json:"companyId" form:"companyId" query:"companyId"`
	Name        string     `json:"name" form:"name" query:"name"`
	Description string     `json:"description" form:"description" query:"description"`
	Address     string     `json:"address" form:"address" query:"address"`
	Email       string     `json:"email" form:"email" query:"email"`
	PhoneNumber string     `json:"phoneNumber" form:"phoneNumber" query:"phoneNumber"`
	CompanyName string     `json:"companyName" form:"companyName" query:"companyName"`
	CreateName  string     `json:"createName" form:"createName" query:"createName"`
	Preloads    string     `json:"preloads" form:"preloads" query:"preloads"`
	StartDt     *time.Time `json:"startDt" form:"startDt" query:"startDt"`
	EndDt       *time.Time `json:"endDt" form:"endDt" query:"endDt"`
}

type CreateCustomer struct {
	CompanyID   string `json:"companyId" form:"companyId" query:"companyId" validate:"required"`
	Name        string `json:"name" form:"name" query:"name" validate:"required"`
	Description string `json:"description" form:"description" query:"description" validate:"required"`
	Address     string `json:"address" form:"address" query:"address" validate:"required"`
	Email       string `json:"email" form:"email" query:"email" validate:"required,email"`
	PhoneNumber string `json:"phoneNumber" form:"phoneNumber" query:"phoneNumber"`
}

type UpdateCustomer struct {
	Name        string `json:"name" form:"name" query:"name" validate:"required"`
	Description string `json:"description" form:"description" query:"description" validate:"required"`
	Address     string `json:"address" form:"address" query:"address" validate:"required"`
	Email       string `json:"email" form:"email" query:"email" validate:"required,email"`
	PhoneNumber string `json:"phoneNumber" form:"phoneNumber" query:"phoneNumber"`
}
