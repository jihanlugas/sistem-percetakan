package request

type CreateUsercompany struct {
	CompanyID        string `json:"companyId" validate:"required"`
	UserID           string `json:"userId" validate:"required"`
	IsDefaultCompany bool   `json:"isDefaultCompany" validate:""`
	IsCreator        bool   `json:"isCreator" validate:""`
}

type UpdateUsercompany struct {
	CompanyID        string `json:"companyId" validate:"required"`
	UserID           string `json:"userId" validate:"required"`
	IsDefaultCompany bool   `json:"isDefaultCompany" validate:""`
	IsCreator        bool   `json:"isCreator" validate:""`
}

type PageUsercompany struct {
	Paging
	CompanyID string `json:"companyId" form:"companyId" query:"companyId"`
	UserID    string `json:"userId" form:"userId" query:"userId"`
}
