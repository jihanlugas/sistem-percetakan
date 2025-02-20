package request

import "time"

type PagePrint struct {
	Paging
	CompanyID       string     `json:"companyId" form:"companyId" query:"companyId"`
	OrderID         string     `json:"orderId" form:"orderId" query:"orderId"`
	PaperID         string     `json:"paperId" form:"paperId" query:"paperId"`
	Name            string     `json:"name" form:"name" query:"name"`
	Description     string     `json:"description" form:"description" query:"description"`
	IsDuplex        *bool      `json:"isDuplex" form:"isDuplex" query:"isDuplex"`
	CompanyName     string     `json:"companyName" form:"companyName" query:"companyName"`
	OrderName       string     `json:"orderName" form:"orderName" query:"orderName"`
	PaperName       string     `json:"paperName" form:"paperName" query:"paperName"`
	CreateName      string     `json:"createName" form:"createName" query:"createName"`
	Preloads        string     `json:"preloads" form:"preloads" query:"preloads"`
	StartDt         *time.Time `json:"startDt" form:"startDt" query:"startDt"`
	EndDt           *time.Time `json:"endDt" form:"endDt" query:"endDt"`
	StartTotalPrint *int64     `json:"startTotalPrint" form:"startTotalPrint" query:"startTotalPrint"`
	EndTotalPrint   *int64     `json:"endTotalPrint" form:"endTotalPrint" query:"endTotalPrint"`
}

type CreatePrint struct {
	CompanyID   string `json:"companyId" form:"companyId" query:"companyId"`
	OrderID     string `json:"orderId" form:"orderId" query:"orderId"`
	Name        string `json:"name" form:"name" query:"name" validate:"required"`
	Description string `json:"description" form:"description" query:"description" validate:""`
	PaperID     string `json:"paperId" form:"paperId" query:"paperId"`
	IsDuplex    bool   `json:"isDuplex" form:"isDuplex" query:"isDuplex"`
	PageCount   int64  `json:"pageCount" form:"pageCount" query:"pageCount" validate:""`
	Qty         int64  `json:"qty" form:"qty" query:"qty" validate:""`
	Price       int64  `json:"price" form:"price" query:"price" validate:""`
	Total       int64  `json:"total" form:"total" query:"total" validate:""`
}

type UpdatePrint struct {
	Name        string `json:"name" form:"name" query:"name" validate:"required"`
	Description string `json:"description" form:"description" query:"description" validate:""`
	PaperID     string `json:"paperId" form:"paperId" query:"paperId"`
	IsDuplex    bool   `json:"isDuplex" form:"isDuplex" query:"isDuplex"`
	PageCount   int64  `json:"pageCount" form:"pageCount" query:"pageCount" validate:""`
	Qty         int64  `json:"qty" form:"qty" query:"qty" validate:""`
	Price       int64  `json:"price" form:"price" query:"price" validate:""`
	Total       int64  `json:"total" form:"total" query:"total" validate:""`
}
