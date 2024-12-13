package request

type PageOrder struct {
	Paging
	CompanyID    string `json:"companyId" form:"companyId" query:"companyId"`
	CustomerID   string `json:"customerId" form:"customerId" query:"customerId"`
	Name         string `json:"name" form:"name" query:"name"`
	Description  string `json:"description" form:"description" query:"description"`
	CompanyName  string `json:"companyName" form:"companyName" query:"companyName"`
	CustomerName string `json:"customerName" form:"customerName" query:"customerName"`
	CreateName   string `json:"createName" form:"createName" query:"createName"`
	Company      bool   `json:"company" form:"company" query:"company"`
	Customer     bool   `json:"customer" form:"customer" query:"customer"`
}
