package request

type CreateCompany struct {
	Fullname    string `json:"fullname" form:"fullname" validate:"required,lte=80"`
	Email       string `json:"email" form:"email" validate:"required,lte=200,email,notexists=email"`
	PhoneNumber string `json:"phoneNumber" form:"phoneNumber" validate:"required,lte=20,notexists=no_hp"`
	Username    string `json:"username" form:"username" validate:"required,lte=20,lowercase,notexists=username"`
	Passwd      string `json:"passwd" form:"passwd" validate:"required,lte=200"`
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:""`
	Balance     int64  `json:"balance" validate:""`
}

type UpdateCompany struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:""`
	Balance     int64  `json:"balance" validate:""`
}

type PageCompany struct {
	Paging
	Name        string `json:"name" form:"name" query:"name"`
	Description string `json:"description" form:"description" query:"description"`
	CreateName  string `json:"createName" form:"createName" query:"createName"`
}
