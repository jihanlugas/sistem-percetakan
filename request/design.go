package request

import "time"

type PageDesign struct {
	Paging
	CompanyID   string     `json:"companyId" form:"companyId" query:"companyId"`
	OrderID     string     `json:"orderId" form:"orderId" query:"orderId"`
	Name        string     `json:"name" form:"name" query:"name"`
	Description string     `json:"description" form:"description" query:"description"`
	CompanyName string     `json:"companyName" form:"companyName" query:"companyName"`
	OrderName   string     `json:"orderName" form:"orderName" query:"orderName"`
	CreateName  string     `json:"createName" form:"createName" query:"createName"`
	Preloads    string     `json:"preloads" form:"preloads" query:"preloads"`
	StartDt     *time.Time `json:"startDt" form:"startDt" query:"startDt"`
	EndDt       *time.Time `json:"endDt" form:"endDt" query:"endDt"`
}

type CreateDesign struct {
	CompanyID   string `json:"companyId" form:"companyId" query:"companyId"`
	OrderID     string `json:"orderId" form:"orderId" query:"orderId"`
	Name        string `json:"name" form:"name" query:"name" validate:"required"`
	Description string `json:"description" form:"description" query:"description" validate:""`
	Qty         int64  `json:"qty" form:"qty" query:"qty" validate:""`
	Price       int64  `json:"price" form:"price" query:"price" validate:""`
	Total       int64  `json:"total" form:"total" query:"total" validate:""`
}

type UpdateDesign struct {
	Name        string `json:"name" form:"name" query:"name" validate:"required"`
	Description string `json:"description" form:"description" query:"description" validate:""`
	Qty         int64  `json:"qty" form:"qty" query:"qty" validate:""`
	Price       int64  `json:"price" form:"price" query:"price" validate:""`
	Total       int64  `json:"total" form:"total" query:"total" validate:""`
}
