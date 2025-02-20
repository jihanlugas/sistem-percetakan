package request

import "time"

type ChangePassword struct {
	CurrentPasswd string `json:"currentPasswd" form:"currentPasswd" validate:"required,lte=200"`
	Passwd        string `json:"passwd" form:"passwd" validate:"required,lte=200"`
	ConfirmPasswd string `json:"confirmPasswd" form:"confirmPasswd" validate:"required,lte=200,eqfield=Passwd"`
}

type CreateUser struct {
	CompanyID   string     `json:"companyId" form:"companyId" query:"companyId"`
	Fullname    string     `json:"fullname" form:"fullname" validate:"required,lte=80"`
	Email       string     `json:"email" form:"email" validate:"required,lte=200,email,notexists=email"`
	PhoneNumber string     `json:"phoneNumber" form:"phoneNumber" validate:"required,lte=20"`
	Username    string     `json:"username" form:"username" validate:"required,lte=20,lowercase,notexists=username"`
	Passwd      string     `json:"passwd" form:"passwd" validate:"required,lte=200"`
	Address     string     `json:"address" form:"address" validate:""`
	BirthDt     *time.Time `json:"birthDt" form:"birthDt" validate:""`
	BirthPlace  string     `json:"birthPlace" form:"birthPlace" validate:""`
}

type UpdateUser struct {
	Fullname    string     `json:"fullname" form:"fullname" validate:"required,lte=80"`
	Email       string     `json:"email" form:"email" validate:"required,lte=200,email"`
	PhoneNumber string     `json:"phoneNumber" form:"phoneNumber" validate:"required,lte=20"`
	Username    string     `json:"username" form:"username" validate:"required,lte=20,lowercase"`
	Address     string     `json:"address" form:"address" validate:""`
	BirthDt     *time.Time `json:"birthDt" form:"birthDt" validate:""`
	BirthPlace  string     `json:"birthPlace" form:"birthPlace" validate:""`
}

type PageUser struct {
	Paging
	CompanyID   string `json:"companyId" form:"companyId" query:"companyId"`
	Fullname    string `json:"fullname" form:"fullname" query:"fullname"`
	Email       string `json:"email" form:"email" query:"email"`
	PhoneNumber string `json:"phoneNumber" form:"phoneNumber" query:"phoneNumber"`
	Username    string `json:"username" form:"username" query:"username"`
	Address     string `json:"address" form:"address" query:"address"`
	BirthPlace  string `json:"birthPlace" form:"birthPlace" query:"birthPlace"`
	CreateName  string `json:"createName" form:"createName" query:"createName"`
	Preloads    string `json:"preloads" form:"preloads" query:"preloads"`
}
