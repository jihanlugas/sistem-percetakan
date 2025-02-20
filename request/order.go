package request

import "time"

type PageOrder struct {
	Paging
	CompanyID       string     `json:"companyId" form:"companyId" query:"companyId"`
	Preloads        string     `json:"preloads" form:"preloads" query:"preloads"`
	CustomerID      string     `json:"customerId" form:"customerId" query:"customerId"`
	PhaseID         string     `json:"phaseId" form:"phaseId" query:"phaseId"`
	Name            string     `json:"name" form:"name" query:"name"`
	Description     string     `json:"description" form:"description" query:"description"`
	IsDone          *bool      `json:"isDone" form:"isDone" query:"isDone"`
	StartDt         *time.Time `json:"startDt" form:"startDt" query:"startDt"`
	EndDt           *time.Time `json:"endDt" form:"endDt" query:"endDt"`
	StartTotalOrder *int64     `json:"startTotalOrder" form:"startTotalOrder" query:"startTotalOrder"`
	EndTotalOrder   *int64     `json:"endTotalOrder" form:"endTotalOrder" query:"endTotalOrder"`
}

type CreateOrder struct {
	CompanyID        string                 `json:"companyId" form:"companyId" query:"companyId" validate:"required"`
	CustomerID       string                 `json:"customerId" form:"customerId" query:"customerId"`
	OrderphaseID     string                 `json:"orderphaseId" form:"orderphaseId" query:"orderphaseId" validate:""`
	NewCustomer      string                 `json:"newCustomer" form:"newCustomer" query:"newCustomer" validate:""`
	NewCustomerPhone string                 `json:"newCustomerPhone" form:"newCustomerPhone" query:"newCustomerPhone" validate:""`
	Name             string                 `json:"name" form:"name" query:"name" validate:"required"`
	Description      string                 `json:"description" form:"description" query:"description" validate:""`
	Designs          []CreateOrderDesign    `json:"designs" form:"designs" query:"designs" validate:"dive"`
	Prints           []CreateOrderPrint     `json:"prints" form:"prints" query:"prints" validate:"dive"`
	Finishings       []CreateOrderFinishing `json:"finishings" form:"finishings" query:"finishings" validate:"dive"`
	Others           []CreateOrderOther     `json:"others" form:"others" query:"others" validate:"dive"`
}

type CreateOrderDesign struct {
	Name        string `json:"name" form:"name" query:"name" validate:"required"`
	Description string `json:"description" form:"description" query:"description" validate:""`
	Qty         int64  `json:"qty" form:"qty" query:"qty" validate:""`
	Price       int64  `json:"price" form:"price" query:"price" validate:""`
	Total       int64  `json:"total" form:"total" query:"total" validate:""`
}

type UpdateOrder struct {
	Name        string `json:"name" form:"name" query:"name" validate:"required"`
	Description string `json:"description" form:"description" query:"description" validate:""`
	IsDone      bool   `json:"isDone" form:"isDone" query:"isDone" validate:""`
}

type CreateOrderPrint struct {
	Name        string `json:"name" form:"name" query:"name" validate:"required"`
	PaperID     string `json:"paperId" form:"paperId" query:"paperId" validate:""`
	Description string `json:"description" form:"description" query:"description" validate:""`
	IsDuplex    bool   `json:"isDuplex" form:"isDuplex" query:"isDuplex" validate:""`
	PageCount   int64  `json:"pageCount" form:"pageCount" query:"pageCount" validate:""`
	Qty         int64  `json:"qty" form:"qty" query:"qty" validate:""`
	Price       int64  `json:"price" form:"price" query:"price" validate:""`
	Total       int64  `json:"total" form:"total" query:"total" validate:""`
}

type CreateOrderFinishing struct {
	Name        string `json:"name" form:"name" query:"name" validate:"required"`
	Description string `json:"description" form:"description" query:"description" validate:""`
	Qty         int64  `json:"qty" form:"qty" query:"qty" validate:""`
	Price       int64  `json:"price" form:"price" query:"price" validate:""`
	Total       int64  `json:"total" form:"total" query:"total" validate:""`
}

type CreateOrderOther struct {
	Name        string `json:"name" form:"name" query:"name" validate:"required"`
	Description string `json:"description" form:"description" query:"description" validate:""`
	Qty         int64  `json:"qty" form:"qty" query:"qty" validate:""`
	Price       int64  `json:"price" form:"price" query:"price" validate:""`
	Total       int64  `json:"total" form:"total" query:"total" validate:""`
}

type AddPhase struct {
	OrderphaseID string `json:"orderphaseId" form:"orderphaseId" query:"orderphaseId" validate:"required"`
}

type AddTransaction struct {
	Name        string `json:"name" form:"name" query:"name" validate:"required"`
	Description string `json:"description" form:"description" query:"description" validate:""`
	Amount      int64  `json:"amount" form:"amount" query:"amount" validate:"required"`
}
