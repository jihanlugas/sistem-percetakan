package request

import "time"

type PageFinishing struct {
	Paging
	CompanyID           string     `json:"companyId" form:"companyId" query:"companyId"`
	OrderID             string     `json:"orderId" form:"orderId" query:"orderId"`
	Name                string     `json:"name" form:"name" query:"name"`
	Description         string     `json:"description" form:"description" query:"description"`
	CompanyName         string     `json:"companyName" form:"companyName" query:"companyName"`
	OrderName           string     `json:"orderName" form:"orderName" query:"orderName"`
	CreateName          string     `json:"createName" form:"createName" query:"createName"`
	Preloads            string     `json:"preloads" form:"preloads" query:"preloads"`
	StartDt             *time.Time `json:"startDt" form:"startDt" query:"startDt"`
	EndDt               *time.Time `json:"endDt" form:"endDt" query:"endDt"`
	StartTotalFinishing *int64     `json:"startTotalFinishing" form:"startTotalFinishing" query:"startTotalFinishing"`
	EndTotalFinishing   *int64     `json:"endTotalFinishing" form:"endTotalFinishing" query:"endTotalFinishing"`
}

type CreateFinishing struct {
	CompanyID   string `json:"companyId" form:"companyId" query:"companyId"`
	OrderID     string `json:"orderId" form:"orderId" query:"orderId"`
	Name        string `json:"name" form:"name" query:"name" validate:"required"`
	Description string `json:"description" form:"description" query:"description" validate:""`
	Qty         int64  `json:"qty" form:"qty" query:"qty" validate:""`
	Price       int64  `json:"price" form:"price" query:"price" validate:""`
	Total       int64  `json:"total" form:"total" query:"total" validate:""`
}

type UpdateFinishing struct {
	Name        string `json:"name" form:"name" query:"name" validate:"required"`
	Description string `json:"description" form:"description" query:"description" validate:""`
	Qty         int64  `json:"qty" form:"qty" query:"qty" validate:""`
	Price       int64  `json:"price" form:"price" query:"price" validate:""`
	Total       int64  `json:"total" form:"total" query:"total" validate:""`
}
