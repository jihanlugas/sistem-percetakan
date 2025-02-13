package request

type PagePayment struct {
	Paging
	CompanyID   string `json:"companyId" form:"companyId" query:"companyId"`
	OrderID     string `json:"orderId" form:"orderId" query:"orderId"`
	Name        string `json:"name" form:"name" query:"name"`
	Description string `json:"description" form:"description" query:"description"`
	CompanyName string `json:"companyName" form:"companyName" query:"companyName"`
	OrderName   string `json:"orderName" form:"orderName" query:"orderName"`
	CreateName  string `json:"createName" form:"createName" query:"createName"`
	Preloads    string `json:"preloads" form:"preloads" query:"preloads"`
	//Company       bool   `json:"company" form:"company" query:"company"`
	//Order         bool   `json:"order" form:"order" query:"order"`
}

type CreatePayment struct {
	CompanyID   string `json:"companyId" form:"companyId" query:"companyId" validate:"required"`
	OrderID     string `json:"orderId" form:"orderId" query:"orderId" validate:"required"`
	Name        string `json:"name" form:"name" query:"name" validate:"required"`
	Description string `json:"description" form:"description" query:"description" validate:""`
	Amount      int64  `json:"amount" form:"amount" query:"amount" validate:"required"`
}

type UpdatePayment struct {
	Name        string `json:"name" form:"name" query:"name" validate:"required"`
	Description string `json:"description" form:"description" query:"description" validate:""`
	Amount      int64  `json:"amount" form:"amount" query:"amount" validate:"required"`
}
