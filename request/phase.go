package request

type PagePhase struct {
	Paging
	CompanyID   string `json:"companyId" form:"companyId" query:"companyId"`
	Name        string `json:"name" form:"name" query:"name"`
	Description string `json:"description" form:"description" query:"description"`
	CompanyName string `json:"companyName" form:"companyName" query:"companyName"`
	CreateName  string `json:"createName" form:"createName" query:"createName"`
	Company     bool   `json:"company" form:"company" query:"company"`
}
