package request

type PageFinishing struct {
	Paging
	CompanyID   string `json:"companyId" form:"companyId" query:"companyId"`
	OrderID     string `json:"orderId" form:"orderId" query:"orderId"`
	Name        string `json:"name" form:"name" query:"name"`
	Description string `json:"description" form:"description" query:"description"`
	CompanyName string `json:"companyName" form:"companyName" query:"companyName"`
	OrderName   string `json:"orderName" form:"orderName" query:"orderName"`
	CreateName  string `json:"createName" form:"createName" query:"createName"`
	Company     bool   `json:"company" form:"company" query:"company"`
	Order       bool   `json:"order" form:"order" query:"order"`
}
