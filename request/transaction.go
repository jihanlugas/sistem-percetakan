package request

import "time"

type PageTransaction struct {
	Paging
	CompanyID   string     `json:"companyId" form:"companyId" query:"companyId"`
	OrderID     string     `json:"orderId" form:"orderId" query:"orderId"`
	Name        string     `json:"name" form:"name" query:"name"`
	Description string     `json:"description" form:"description" query:"description"`
	Type        int64      `json:"type" form:"type" query:"type"`
	Amount      int64      `json:"amount" form:"amount" query:"amount"`
	CompanyName string     `json:"companyName" form:"companyName" query:"companyName"`
	OrderName   string     `json:"orderName" form:"orderName" query:"orderName"`
	CreateName  string     `json:"createName" form:"createName" query:"createName"`
	Preloads    string     `json:"preloads" form:"preloads" query:"preloads"`
	StartDt     *time.Time `json:"startDt" form:"startDt" query:"startDt"`
	EndDt       *time.Time `json:"endDt" form:"endDt" query:"endDt"`
	StartAmount *int64     `json:"startAmount" form:"startAmount" query:"startAmount"`
	EndAmount   *int64     `json:"endAmount" form:"endAmount" query:"endAmount"`
}

type CreateTransaction struct {
	CompanyID   string `json:"companyId" form:"companyId" query:"companyId" validate:"required"`
	OrderID     string `json:"orderId" form:"orderId" query:"orderId" validate:""`
	Name        string `json:"name" form:"name" query:"name" validate:"required"`
	Description string `json:"description" form:"description" query:"description" validate:""`
	Type        int64  `json:"type" form:"type" query:"type"`
	Amount      int64  `json:"amount" form:"amount" query:"amount" validate:"required"`
}

type UpdateTransaction struct {
	Name        string `json:"name" form:"name" query:"name" validate:"required"`
	Description string `json:"description" form:"description" query:"description" validate:""`
	Type        int64  `json:"type" form:"type" query:"type"`
	Amount      int64  `json:"amount" form:"amount" query:"amount" validate:"required"`
}
