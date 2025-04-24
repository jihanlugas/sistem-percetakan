package model

import "time"

const (
	VIEW_PHOTO       = "photos_view"
	VIEW_USER        = "users_view"
	VIEW_COMPANY     = "companies_view"
	VIEW_USERCOMPANY = "usercompanies_view"
	VIEW_ORDER       = "orders_view"
	VIEW_PRINT       = "prints_view"
	VIEW_FINISHING   = "finishings_view"
	VIEW_CUSTOMER    = "customers_view"
	VIEW_PAPER       = "papers_view"
	VIEW_PHASE       = "phases_view"
	VIEW_ORDERPHASE  = "orderphases_view"
	VIEW_TRANSACTION = "transactions_view"
)

type UserLogin struct {
	ExpiredDt     time.Time `json:"expiredDt"`
	UserID        string    `json:"userId"`
	PassVersion   int       `json:"passVersion"`
	CompanyID     string    `json:"companyId"`
	Role          string    `json:"role"`
	UsercompanyID string    `json:"usercompanyId"`
}
