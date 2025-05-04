package request

type UpdateCompany struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:""`
	Email       string `json:"email" validate:"required"`
	PhoneNumber string `json:"phoneNumber" validate:"required"`
	Address     string `json:"address" validate:"required"`
	InvoiceNote string `json:"invoiceNote" validate:""`
}

type PageCompany struct {
	Paging
	Name        string `json:"name" form:"name" query:"name"`
	Description string `json:"description" form:"description" query:"description"`
	CreateName  string `json:"createName" form:"createName" query:"createName"`
}
