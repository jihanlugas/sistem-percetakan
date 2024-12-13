package request

type ChangePassword struct {
	CurrentPasswd string `json:"currentPasswd" form:"currentPasswd" validate:"required,lte=200"`
	Passwd        string `json:"passwd" form:"passwd" validate:"required,lte=200"`
	ConfirmPasswd string `json:"confirmPasswd" form:"confirmPasswd" validate:"required,lte=200,eqfield=Passwd"`
}

type CreateUser struct {
	Fullname string `json:"fullname" form:"fullname" validate:"required,lte=80"`
	Email    string `json:"email" form:"email" validate:"required,lte=200,email,notexists=email"`
	NoHp     string `json:"noHp" form:"noHp" validate:"required,lte=20,notexists=no_hp"`
	Username string `json:"username" form:"username" validate:"required,lte=20,lowercase,notexists=username"`
	Passwd   string `json:"passwd" form:"passwd" validate:"required,lte=200"`
}

type UpdateUser struct {
	Fullname string `json:"fullname" form:"fullname" validate:"required,lte=80"`
	Email    string `json:"email" form:"email" validate:"required,lte=200,email,notexists=email"`
	NoHp     string `json:"noHp" form:"noHp" validate:"required,lte=20,notexists=no_hp"`
	Username string `json:"username" form:"username" validate:"required,lte=20,lowercase,notexists=username"`
}

type PageUser struct {
	Paging
	Email    string `json:"email" form:"email" query:"email"`
	NoHp     string `json:"noHp" form:"noHp" query:"noHp"`
	Username string `json:"username" form:"username" query:"username"`
	Fullname string `json:"fullname" form:"fullname" query:"fullname"`
}
