package request

type PagePaper struct {
	Paging
	CompanyID   string `json:"companyId" form:"companyId" query:"companyId"`
	Name        string `json:"name" form:"name" query:"name"`
	Description string `json:"description" form:"description" query:"description"`
	CompanyName string `json:"companyName" form:"companyName" query:"companyName"`
	CreateName  string `json:"createName" form:"createName" query:"createName"`
	Preloads    string `json:"preloads" form:"preloads" query:"preloads"`
}

type CreatePaper struct {
	CompanyID          string `json:"companyId" form:"companyId" query:"companyId"`
	Name               string `json:"name" form:"name" query:"name" validate:"required"`
	Description        string `json:"description" form:"description" query:"description" validate:""`
	DefaultPrice       int64  `json:"defaultPrice" form:"defaultPrice" query:"defaultPrice"`
	DefaultPriceDuplex int64  `json:"defaultPriceDuplex" form:"defaultPriceDuplex" query:"defaultPriceDuplex"`
}

type UpdatePaper struct {
	Name               string `json:"name" form:"name" query:"name" validate:"required"`
	Description        string `json:"description" form:"description" query:"description" validate:""`
	DefaultPrice       int64  `json:"defaultPrice" form:"defaultPrice" query:"defaultPrice"`
	DefaultPriceDuplex int64  `json:"defaultPriceDuplex" form:"defaultPriceDuplex" query:"defaultPriceDuplex"`
}
