package model

import (
	"gorm.io/gorm"
	"time"
)

type PhotoView struct {
	ID          string         `json:"id"`
	ClientName  string         `json:"clientName"`
	ServerName  string         `json:"serverName"`
	RefTable    string         `json:"refTable"`
	Ext         string         `json:"ext"`
	PhotoPath   string         `json:"photoPath"`
	PhotoSize   int64          `json:"photoSize"`
	PhotoWidth  int64          `json:"photoWidth"`
	PhotoHeight int64          `json:"photoHeight"`
	CreateBy    string         `json:"createBy"`
	CreateDt    time.Time      `json:"createDt"`
	DeleteDt    gorm.DeletedAt `json:"deleteDt"`
	CreateName  string         `json:"createName"`
}

func (PhotoView) TableName() string {
	return VIEW_PHOTO
}

type UserView struct {
	ID                string         `json:"id"`
	Role              string         `json:"role"`
	Email             string         `json:"email"`
	Username          string         `json:"username"`
	NoHp              string         `json:"noHp"`
	Fullname          string         `json:"fullname"`
	Passwd            string         `json:"-"`
	PassVersion       int            `json:"passVersion"`
	IsActive          bool           `json:"isActive"`
	PhotoID           string         `json:"photoId"`
	PhotoUrl          string         `json:"photoUrl"`
	LastLoginDt       *time.Time     `json:"lastLoginDt"`
	BirthDt           *time.Time     `json:"birthDt"`
	BirthPlace        string         `json:"birthPlace"`
	AccountVerifiedDt *time.Time     `json:"accountVerifiedDt"`
	CreateBy          string         `json:"createBy"`
	CreateDt          time.Time      `json:"createDt"`
	UpdateBy          string         `json:"updateBy"`
	UpdateDt          time.Time      `json:"updateDt"`
	DeleteDt          gorm.DeletedAt `json:"deleteDt"`
	CreateName        string         `json:"createName"`
	UpdateName        string         `json:"updateName"`
}

func (UserView) TableName() string {
	return VIEW_USER
}

type CompanyView struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	PhotoID     string         `json:"photoId"`
	PhotoUrl    string         `json:"photoUrl"`
	CreateBy    string         `json:"createBy"`
	CreateDt    time.Time      `json:"createDt"`
	UpdateBy    string         `json:"updateBy"`
	UpdateDt    time.Time      `json:"updateDt"`
	DeleteDt    gorm.DeletedAt `json:"deleteDt"`
	CreateName  string         `json:"createName"`
	UpdateName  string         `json:"updateName"`
}

func (CompanyView) TableName() string {
	return VIEW_COMPANY
}

type UsercompanyView struct {
	ID               string         `json:"id"`
	UserID           string         `json:"userId"`
	CompanyID        string         `json:"companyId"`
	IsDefaultCompany bool           `json:"isDefaultCompany"`
	IsCreator        bool           `json:"isCreator"`
	CreateBy         string         `json:"createBy"`
	CreateDt         time.Time      `json:"createDt"`
	UpdateBy         string         `json:"updateBy"`
	UpdateDt         time.Time      `json:"updateDt"`
	DeleteDt         gorm.DeletedAt `json:"deleteDt"`
	UserName         string         `json:"userName"`
	CompanyName      string         `json:"companyName"`
	CreateName       string         `json:"createName"`
	UpdateName       string         `json:"updateName"`
}

func (UsercompanyView) TableName() string {
	return VIEW_USERCOMPANY
}

type CustomerView struct {
	ID          string         `json:"id"`
	CompanyID   string         `json:"companyId"`
	Name        string         `json:"name"`
	Customer    string         `json:"customer"`
	NoHp        string         `json:"noHp"`
	CreateBy    string         `json:"createBy"`
	CreateDt    time.Time      `json:"createDt"`
	UpdateBy    string         `json:"updateBy"`
	UpdateDt    time.Time      `json:"updateDt"`
	DeleteDt    gorm.DeletedAt `json:"deleteDt"`
	CompanyName string         `json:"companyName"`
	CreateName  string         `json:"createName"`
	UpdateName  string         `json:"updateName"`
}

func (CustomerView) TableName() string {
	return VIEW_CUSTOMER
}

type OrderView struct {
	ID           string         `json:"id"`
	CompanyID    string         `json:"companyId"`
	CustomerID   string         `json:"customerId"`
	Name         string         `json:"name"`
	Description  string         `json:"description"`
	IsDone       bool           `json:"isDone"`
	CreateBy     string         `json:"createBy"`
	CreateDt     time.Time      `json:"createDt"`
	UpdateBy     string         `json:"updateBy"`
	UpdateDt     time.Time      `json:"updateDt"`
	DeleteDt     gorm.DeletedAt `json:"deleteDt"`
	CompanyName  string         `json:"companyName"`
	CustomerName string         `json:"customerName"`
	CreateName   string         `json:"createName"`
	UpdateName   string         `json:"updateName"`

	Company  *CompanyView  `json:"company"`
	Customer *CustomerView `json:"customer"`
}

func (OrderView) TableName() string {
	return VIEW_ORDER
}

type PaperView struct {
	ID          string         `json:"id"`
	CompanyID   string         `json:"companyId"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	CreateBy    string         `json:"createBy"`
	CreateDt    time.Time      `json:"createDt"`
	UpdateBy    string         `json:"updateBy"`
	UpdateDt    time.Time      `json:"updateDt"`
	DeleteDt    gorm.DeletedAt `json:"deleteDt"`
	CompanyName string         `json:"companyName"`
	CreateName  string         `json:"createName"`
	UpdateName  string         `json:"updateName"`
}

func (PaperView) TableName() string {
	return VIEW_PAPER
}

type DesignView struct {
	ID          string         `json:"id"`
	CompanyID   string         `json:"companyId"`
	OrderID     string         `json:"orderId"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Qty         int64          `json:"qty"`
	Price       int64          `json:"price"`
	Total       int64          `json:"total"`
	CreateBy    string         `json:"createBy"`
	CreateDt    time.Time      `json:"createDt"`
	UpdateBy    string         `json:"updateBy"`
	UpdateDt    time.Time      `json:"updateDt"`
	DeleteDt    gorm.DeletedAt `json:"deleteDt"`
	CompanyName string         `json:"companyName"`
	OrderName   string         `json:"orderName"`
	CreateName  string         `json:"createName"`
	UpdateName  string         `json:"updateName"`
}

func (DesignView) TableName() string {
	return VIEW_DESIGN
}

type PrintView struct {
	ID          string         `json:"id"`
	CompanyID   string         `json:"companyId"`
	OrderID     string         `json:"orderId"`
	PaperID     string         `json:"paperId"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	IsDuplex    bool           `json:"isDuplex"`
	PageCount   int64          `json:"pageCount"`
	PrintCount  int64          `json:"printCount"`
	Price       int64          `json:"price"`
	Total       int64          `json:"total"`
	CreateBy    string         `json:"createBy"`
	CreateDt    time.Time      `json:"createDt"`
	UpdateBy    string         `json:"updateBy"`
	UpdateDt    time.Time      `json:"updateDt"`
	DeleteDt    gorm.DeletedAt `json:"deleteDt"`
	CompanyName string         `json:"companyName"`
	OrderName   string         `json:"orderName"`
	CreateName  string         `json:"createName"`
	UpdateName  string         `json:"updateName"`
}

func (PrintView) TableName() string {
	return VIEW_PRINT
}

type FinishingView struct {
	ID          string         `json:"id"`
	CompanyID   string         `json:"companyId"`
	OrderID     string         `json:"orderId"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Qty         int64          `json:"qty"`
	Price       int64          `json:"price"`
	Total       int64          `json:"total"`
	CreateBy    string         `json:"createBy"`
	CreateDt    time.Time      `json:"createDt"`
	UpdateBy    string         `json:"updateBy"`
	UpdateDt    time.Time      `json:"updateDt"`
	DeleteDt    gorm.DeletedAt `json:"deleteDt"`
	CompanyName string         `json:"companyName"`
	OrderName   string         `json:"orderName"`
	CreateName  string         `json:"createName"`
	UpdateName  string         `json:"updateName"`
}

func (FinishingView) TableName() string {
	return VIEW_FINISHING
}

type OtherView struct {
	ID          string         `json:"id"`
	CompanyID   string         `json:"companyId"`
	OrderID     string         `json:"orderId"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Qty         int64          `json:"qty"`
	Price       int64          `json:"price"`
	Total       int64          `json:"total"`
	CreateBy    string         `json:"createBy"`
	CreateDt    time.Time      `json:"createDt"`
	UpdateBy    string         `json:"updateBy"`
	UpdateDt    time.Time      `json:"updateDt"`
	DeleteDt    gorm.DeletedAt `json:"deleteDt"`
	CompanyName string         `json:"companyName"`
	OrderName   string         `json:"orderName"`
	CreateName  string         `json:"createName"`
	UpdateName  string         `json:"updateName"`
}

func (OtherView) TableName() string {
	return VIEW_OTHER
}

type PhaseView struct {
	ID          string         `json:"id"`
	CompanyID   string         `json:"companyId"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Order       int64          `json:"order"`
	CreateBy    string         `json:"createBy"`
	CreateDt    time.Time      `json:"createDt"`
	UpdateBy    string         `json:"updateBy"`
	UpdateDt    time.Time      `json:"updateDt"`
	DeleteDt    gorm.DeletedAt `json:"deleteDt"`
	CompanyName string         `json:"companyName"`
	CreateName  string         `json:"createName"`
	UpdateName  string         `json:"updateName"`
}

func (PhaseView) TableName() string {
	return VIEW_PHASE
}

type OrderphaseView struct {
	ID          string         `json:"id"`
	CompanyID   string         `json:"companyId"`
	OrderID     string         `json:"orderId"`
	Name        string         `json:"name"`
	CreateBy    string         `json:"createBy"`
	CreateDt    time.Time      `json:"createDt"`
	UpdateBy    string         `json:"updateBy"`
	UpdateDt    time.Time      `json:"updateDt"`
	DeleteDt    gorm.DeletedAt `json:"deleteDt"`
	CompanyName string         `json:"companyName"`
	OrderName   string         `json:"orderName"`
	CreateName  string         `json:"createName"`
	UpdateName  string         `json:"updateName"`
}

func (OrderphaseView) TableName() string {
	return VIEW_ORDERPHASE
}

type PaymentView struct {
	ID            string         `json:"id"`
	CompanyID     string         `json:"companyId"`
	OrderID       string         `json:"orderId"`
	Name          string         `json:"name"`
	Description   string         `json:"description"`
	IsDonePayment bool           `json:"isDonePayment"`
	Amount        int64          `json:"amount"`
	CreateBy      string         `json:"createBy"`
	CreateDt      time.Time      `json:"createDt"`
	UpdateBy      string         `json:"updateBy"`
	UpdateDt      time.Time      `json:"updateDt"`
	DeleteDt      gorm.DeletedAt `json:"deleteDt"`
	CompanyName   string         `json:"companyName"`
	OrderName     string         `json:"orderName"`
	CreateName    string         `json:"createName"`
	UpdateName    string         `json:"updateName"`
}

func (PaymentView) TableName() string {
	return VIEW_PAYMENT
}
