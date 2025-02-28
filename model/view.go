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
	CompanyID         string         `json:"companyId"`
	UsercompanyID     string         `json:"usercompanyId"`
	Role              string         `json:"role"`
	Email             string         `json:"email"`
	Username          string         `json:"username"`
	PhoneNumber       string         `json:"phoneNumber"`
	Address           string         `json:"address"`
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

	Company *CompanyView `json:"company,omitempty"`
}

func (UserView) TableName() string {
	return VIEW_USER
}

type CompanyView struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Email       string         `json:"email"`
	PhoneNumber string         `json:"phoneNumber"`
	Address     string         `json:"address"`
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

	User    *UserView    `json:"user,omitempty"`
	Company *CompanyView `json:"company,omitempty"`
}

func (UsercompanyView) TableName() string {
	return VIEW_USERCOMPANY
}

type CustomerView struct {
	ID          string         `json:"id"`
	CompanyID   string         `json:"companyId"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Email       string         `json:"email"`
	Address     string         `json:"address"`
	PhoneNumber string         `json:"phoneNumber"`
	CreateBy    string         `json:"createBy"`
	CreateDt    time.Time      `json:"createDt"`
	UpdateBy    string         `json:"updateBy"`
	UpdateDt    time.Time      `json:"updateDt"`
	DeleteDt    gorm.DeletedAt `json:"deleteDt"`
	CompanyName string         `json:"companyName"`
	CreateName  string         `json:"createName"`
	UpdateName  string         `json:"updateName"`

	Company *CompanyView `json:"company,omitempty"`
}

func (CustomerView) TableName() string {
	return VIEW_CUSTOMER
}

type OrderView struct {
	ID               string         `json:"id"`
	CompanyID        string         `json:"companyId"`
	CustomerID       string         `json:"customerId"`
	Name             string         `json:"name"`
	Description      string         `json:"description"`
	Number           int64          `json:"number"`
	CreateBy         string         `json:"createBy"`
	CreateDt         time.Time      `json:"createDt"`
	UpdateBy         string         `json:"updateBy"`
	UpdateDt         time.Time      `json:"updateDt"`
	DeleteDt         gorm.DeletedAt `json:"deleteDt"`
	OrderphaseID     string         `json:"orderphaseId"`
	PhaseID          string         `json:"phaseId"`
	OrderphaseName   string         `json:"orderphaseName"`
	TotalDesign      int64          `json:"totalDesign"`
	TotalPrint       int64          `json:"totalPrint"`
	TotalFinishing   int64          `json:"totalFinishing"`
	TotalOther       int64          `json:"totalOther"`
	TotalTransaction int64          `json:"totalTransaction"`
	TotalOrder       int64          `json:"totalOrder"`
	Outstanding      int64          `json:"outstanding"`
	CompanyName      string         `json:"companyName"`
	CustomerName     string         `json:"customerName"`
	CreateName       string         `json:"createName"`
	UpdateName       string         `json:"updateName"`

	Company      *CompanyView      `json:"company,omitempty"`
	Customer     *CustomerView     `json:"customer,omitempty"`
	Designs      []DesignView      `json:"designs,omitempty" gorm:"foreignKey:OrderID"`
	Prints       []PrintView       `json:"prints,omitempty" gorm:"foreignKey:OrderID"`
	Finishings   []FinishingView   `json:"finishings,omitempty" gorm:"foreignKey:OrderID"`
	Others       []OtherView       `json:"others,omitempty" gorm:"foreignKey:OrderID"`
	Orderphases  []OrderphaseView  `json:"orderphases,omitempty" gorm:"foreignKey:OrderID"`
	Transactions []TransactionView `json:"transactions,omitempty" gorm:"foreignKey:OrderID"`
}

func (OrderView) TableName() string {
	return VIEW_ORDER
}

type PaperView struct {
	ID                 string         `json:"id"`
	CompanyID          string         `json:"companyId"`
	Name               string         `json:"name"`
	Description        string         `json:"description"`
	DefaultPrice       int64          `json:"defaultPrice"`
	DefaultPriceDuplex int64          `json:"defaultPriceDuplex"`
	CreateBy           string         `json:"createBy"`
	CreateDt           time.Time      `json:"createDt"`
	UpdateBy           string         `json:"updateBy"`
	UpdateDt           time.Time      `json:"updateDt"`
	DeleteDt           gorm.DeletedAt `json:"deleteDt"`
	CompanyName        string         `json:"companyName"`
	CreateName         string         `json:"createName"`
	UpdateName         string         `json:"updateName"`

	Company *CompanyView `json:"company,omitempty"`
}

func (PaperView) TableName() string {
	return VIEW_PAPER
}

type DesignView struct {
	ID            string         `json:"id"`
	UsercompanyID string         `json:"usercompanyId"`
	CompanyID     string         `json:"companyId"`
	OrderID       string         `json:"orderId"`
	Name          string         `json:"name"`
	Description   string         `json:"description"`
	Qty           int64          `json:"qty"`
	Price         int64          `json:"price"`
	Total         int64          `json:"total"`
	CreateBy      string         `json:"createBy"`
	CreateDt      time.Time      `json:"createDt"`
	UpdateBy      string         `json:"updateBy"`
	UpdateDt      time.Time      `json:"updateDt"`
	DeleteDt      gorm.DeletedAt `json:"deleteDt"`
	CompanyName   string         `json:"companyName"`
	OrderName     string         `json:"orderName"`
	CreateName    string         `json:"createName"`
	UpdateName    string         `json:"updateName"`

	Company *CompanyView `json:"company,omitempty"`
	Order   *OrderView   `json:"order,omitempty"`
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
	PaperName   string         `json:"paperName"`
	CreateName  string         `json:"createName"`
	UpdateName  string         `json:"updateName"`

	Company *CompanyView `json:"company,omitempty"`
	Order   *OrderView   `json:"order,omitempty"`
	Paper   *PaperView   `json:"paper,omitempty"`
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

	Company *CompanyView `json:"company,omitempty"`
	Order   *OrderView   `json:"order,omitempty"`
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

	Company *CompanyView `json:"company,omitempty"`
	Order   *OrderView   `json:"order,omitempty"`
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

	Company *CompanyView `json:"company,omitempty"`
}

func (PhaseView) TableName() string {
	return VIEW_PHASE
}

type OrderphaseView struct {
	ID          string         `json:"id"`
	CompanyID   string         `json:"companyId"`
	OrderID     string         `json:"orderId"`
	PhaseID     string         `json:"phaseId"`
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

	Company *CompanyView `json:"company,omitempty"`
	Order   *OrderView   `json:"order,omitempty"`
}

func (OrderphaseView) TableName() string {
	return VIEW_ORDERPHASE
}

type TransactionView struct {
	ID          string         `json:"id"`
	CompanyID   string         `json:"companyId"`
	OrderID     string         `json:"orderId"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Type        int64          `json:"type"`
	Amount      int64          `json:"amount"`
	CreateBy    string         `json:"createBy"`
	CreateDt    time.Time      `json:"createDt"`
	UpdateBy    string         `json:"updateBy"`
	UpdateDt    time.Time      `json:"updateDt"`
	DeleteDt    gorm.DeletedAt `json:"deleteDt"`
	CompanyName string         `json:"companyName"`
	OrderName   string         `json:"orderName"`
	CreateName  string         `json:"createName"`
	UpdateName  string         `json:"updateName"`

	Company *CompanyView `json:"company,omitempty"`
	Order   *OrderView   `json:"order,omitempty"`
}

func (TransactionView) TableName() string {
	return VIEW_TRANSACTION
}
