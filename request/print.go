package request

type PagePrint struct {
	Paging
	CompanyID   string `json:"companyId" form:"companyId" query:"companyId"`
	OrderID     string `json:"orderId" form:"orderId" query:"orderId"`
	PaperID     string `json:"paperId" form:"paperId" query:"paperId"`
	Name        string `json:"name" form:"name" query:"name"`
	Description string `json:"description" form:"description" query:"description"`
	IsDuplex    *bool  `json:"isDuplex" form:"isDuplex" query:"isDuplex"`
	CompanyName string `json:"companyName" form:"companyName" query:"companyName"`
	OrderName   string `json:"orderName" form:"orderName" query:"orderName"`
	PaperName   string `json:"paperName" form:"paperName" query:"paperName"`
	CreateName  string `json:"createName" form:"createName" query:"createName"`
	Company     bool   `json:"company" form:"company" query:"company"`
	Order       bool   `json:"order" form:"order" query:"order"`
	Paper       bool   `json:"paper" form:"paper" query:"paper"`
}
