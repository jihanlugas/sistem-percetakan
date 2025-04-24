package model

import (
	"github.com/jihanlugas/sistem-percetakan/constant"
	"gorm.io/gorm"
	"time"
)

type Photo struct {
	ID          string         `gorm:"primaryKey" json:"id"`
	ClientName  string         `gorm:"not null" json:"clientName"`
	ServerName  string         `gorm:"not null" json:"serverName"`
	RefTable    string         `gorm:"not null" json:"refTable"`
	Ext         string         `gorm:"not null" json:"ext"`
	PhotoPath   string         `gorm:"not null" json:"photoPath"`
	PhotoSize   int64          `gorm:"not null" json:"photoSize"`
	PhotoWidth  int64          `gorm:"not null" json:"photoWidth"`
	PhotoHeight int64          `gorm:"not null" json:"photoHeight"`
	CreateBy    string         `gorm:"not null" json:"createBy"`
	CreateDt    time.Time      `gorm:"not null" json:"createDt"`
	DeleteDt    gorm.DeletedAt `gorm:"null" json:"deleteDt"`
}

type Photoinc struct {
	ID        string `gorm:"primaryKey" json:"id"`
	RefTable  string `gorm:"not null" json:"refTable"`
	FolderInc int64  `gorm:"not null" json:"folderInc"`
	Folder    string `gorm:"not null" json:"folder"`
	Running   int64  `gorm:"not null" json:"running"`
}

type User struct {
	ID                string         `gorm:"primaryKey" json:"id"`
	Role              string         `gorm:"not null" json:"role"`
	Email             string         `gorm:"not null" json:"email"`
	Username          string         `gorm:"not null" json:"username"`
	PhoneNumber       string         `gorm:"not null" json:"phoneNumber"`
	Address           string         `gorm:"not null" json:"address"`
	Fullname          string         `gorm:"not null" json:"fullname"`
	Passwd            string         `gorm:"not null" json:"-"`
	PassVersion       int            `gorm:"not null" json:"passVersion"`
	IsActive          bool           `gorm:"not null" json:"isActive"`
	PhotoID           string         `gorm:"not null" json:"photoId"`
	LastLoginDt       *time.Time     `gorm:"null" json:"lastLoginDt"`
	BirthDt           *time.Time     `gorm:"null" json:"birthDt"`
	BirthPlace        string         `gorm:"not null" json:"birthPlace"`
	AccountVerifiedDt *time.Time     `gorm:"null" json:"accountVerifiedDt"`
	CreateBy          string         `gorm:"not null" json:"createBy"`
	CreateDt          time.Time      `gorm:"not null" json:"createDt"`
	UpdateBy          string         `gorm:"not null" json:"updateBy"`
	UpdateDt          time.Time      `gorm:"not null" json:"updateDt"`
	DeleteDt          gorm.DeletedAt `gorm:"null" json:"deleteDt"`
}

type Company struct {
	ID          string         `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"not null" json:"name"`
	Description string         `gorm:"not null" json:"description"`
	Email       string         `gorm:"not null" json:"email"`
	PhoneNumber string         `gorm:"not null" json:"phoneNumber"`
	Address     string         `gorm:"not null" json:"address"`
	PhotoID     string         `gorm:"not null" json:"photoId"`
	CreateBy    string         `gorm:"not null" json:"createBy"`
	CreateDt    time.Time      `gorm:"not null" json:"createDt"`
	UpdateBy    string         `gorm:"not null" json:"updateBy"`
	UpdateDt    time.Time      `gorm:"not null" json:"updateDt"`
	DeleteDt    gorm.DeletedAt `gorm:"null" json:"deleteDt"`
}

type Usercompany struct {
	ID               string         `gorm:"primaryKey" json:"id"`
	UserID           string         `gorm:"not null" json:"userId"`
	CompanyID        string         `gorm:"not null" json:"companyId"`
	IsDefaultCompany bool           `gorm:"not null" json:"isDefaultCompany"`
	IsCreator        bool           `gorm:"not null" json:"isCreator"`
	CreateBy         string         `gorm:"not null" json:"createBy"`
	CreateDt         time.Time      `gorm:"not null" json:"createDt"`
	UpdateBy         string         `gorm:"not null" json:"updateBy"`
	UpdateDt         time.Time      `gorm:"not null" json:"updateDt"`
	DeleteDt         gorm.DeletedAt `gorm:"null" json:"deleteDt"`
}

type Customer struct {
	ID          string         `gorm:"primaryKey" json:"id"`
	CompanyID   string         `gorm:"not null" json:"companyId"`
	Name        string         `gorm:"not null" json:"name"`
	Description string         `gorm:"not null" json:"description"`
	Email       string         `gorm:"not null" json:"email"`
	PhoneNumber string         `gorm:"not null" json:"phoneNumber"`
	Address     string         `gorm:"not null" json:"address"`
	CreateBy    string         `gorm:"not null" json:"createBy"`
	CreateDt    time.Time      `gorm:"not null" json:"createDt"`
	UpdateBy    string         `gorm:"not null" json:"updateBy"`
	UpdateDt    time.Time      `gorm:"not null" json:"updateDt"`
	DeleteDt    gorm.DeletedAt `gorm:"null" json:"deleteDt"`
}

type Order struct {
	ID          string         `gorm:"primaryKey" json:"id"`
	CompanyID   string         `gorm:"not null" json:"companyId"`
	CustomerID  string         `gorm:"not null" json:"customerId"`
	Name        string         `gorm:"not null" json:"name"`
	Description string         `gorm:"not null" json:"description"`
	Number      int64          `gorm:"not null" json:"number"`
	CreateBy    string         `gorm:"not null" json:"createBy"`
	CreateDt    time.Time      `gorm:"not null" json:"createDt"`
	UpdateBy    string         `gorm:"not null" json:"updateBy"`
	UpdateDt    time.Time      `gorm:"not null" json:"updateDt"`
	DeleteDt    gorm.DeletedAt `gorm:"null" json:"deleteDt"`
}

type Paper struct {
	ID                 string         `gorm:"primaryKey" json:"id"`
	CompanyID          string         `gorm:"not null" json:"companyId"`
	Name               string         `gorm:"not null" json:"name"`
	Description        string         `gorm:"not null" json:"description"`
	DefaultPrice       int64          `gorm:"not null" json:"defaultPrice"`
	DefaultPriceDuplex int64          `gorm:"not null" json:"defaultPriceDuplex"`
	CreateBy           string         `gorm:"not null" json:"createBy"`
	CreateDt           time.Time      `gorm:"not null" json:"createDt"`
	UpdateBy           string         `gorm:"not null" json:"updateBy"`
	UpdateDt           time.Time      `gorm:"not null" json:"updateDt"`
	DeleteDt           gorm.DeletedAt `gorm:"null" json:"deleteDt"`
}

type Print struct {
	ID          string         `gorm:"primaryKey" json:"id"`
	CompanyID   string         `gorm:"not null" json:"companyId"`
	OrderID     string         `gorm:"not null" json:"orderId"`
	PaperID     string         `gorm:"not null" json:"paperId"`
	Name        string         `gorm:"not null" json:"name"`
	Description string         `gorm:"not null" json:"description"`
	IsDuplex    bool           `gorm:"not null" json:"isDuplex"`
	PageCount   int64          `gorm:"not null" json:"pageCount"`
	Qty         int64          `gorm:"not null" json:"qty"`
	Price       int64          `gorm:"not null" json:"price"`
	Total       int64          `gorm:"not null" json:"total"`
	CreateBy    string         `gorm:"not null" json:"createBy"`
	CreateDt    time.Time      `gorm:"not null" json:"createDt"`
	UpdateBy    string         `gorm:"not null" json:"updateBy"`
	UpdateDt    time.Time      `gorm:"not null" json:"updateDt"`
	DeleteDt    gorm.DeletedAt `gorm:"null" json:"deleteDt"`
}

type Finishing struct {
	ID          string         `gorm:"primaryKey" json:"id"`
	CompanyID   string         `gorm:"not null" json:"companyId"`
	OrderID     string         `gorm:"not null" json:"orderId"`
	Name        string         `gorm:"not null" json:"name"`
	Description string         `gorm:"not null" json:"description"`
	Qty         int64          `gorm:"not null" json:"qty"`
	Price       int64          `gorm:"not null" json:"price"`
	Total       int64          `gorm:"not null" json:"total"`
	CreateBy    string         `gorm:"not null" json:"createBy"`
	CreateDt    time.Time      `gorm:"not null" json:"createDt"`
	UpdateBy    string         `gorm:"not null" json:"updateBy"`
	UpdateDt    time.Time      `gorm:"not null" json:"updateDt"`
	DeleteDt    gorm.DeletedAt `gorm:"null" json:"deleteDt"`
}

type Phase struct {
	ID          string         `gorm:"primaryKey" json:"id"`
	CompanyID   string         `gorm:"not null" json:"companyId"`
	Name        string         `gorm:"not null" json:"name"`
	Description string         `gorm:"not null" json:"description"`
	Order       int64          `gorm:"not null" json:"order"`
	CreateBy    string         `gorm:"not null" json:"createBy"`
	CreateDt    time.Time      `gorm:"not null" json:"createDt"`
	UpdateBy    string         `gorm:"not null" json:"updateBy"`
	UpdateDt    time.Time      `gorm:"not null" json:"updateDt"`
	DeleteDt    gorm.DeletedAt `gorm:"null" json:"deleteDt"`
}

type Orderphase struct {
	ID        string         `gorm:"primaryKey" json:"id"`
	CompanyID string         `gorm:"not null" json:"companyId"`
	OrderID   string         `gorm:"not null" json:"orderId"`
	PhaseID   string         `gorm:"not null" json:"phaseId"`
	Name      string         `gorm:"not null" json:"name"`
	CreateBy  string         `gorm:"not null" json:"createBy"`
	CreateDt  time.Time      `gorm:"not null" json:"createDt"`
	UpdateBy  string         `gorm:"not null" json:"updateBy"`
	UpdateDt  time.Time      `gorm:"not null" json:"updateDt"`
	DeleteDt  gorm.DeletedAt `gorm:"null" json:"deleteDt"`
}

type Transaction struct {
	ID          string                   `gorm:"primaryKey" json:"id"`
	CompanyID   string                   `gorm:"not null" json:"companyId"`
	OrderID     string                   `gorm:"not null" json:"orderId"`
	Name        string                   `gorm:"not null" json:"name"`
	Description string                   `gorm:"not null" json:"description"`
	Type        constant.TransactionType `gorm:"not null" json:"type"`
	Amount      int64                    `gorm:"not null" json:"amount"`
	CreateBy    string                   `gorm:"not null" json:"createBy"`
	CreateDt    time.Time                `gorm:"not null" json:"createDt"`
	UpdateBy    string                   `gorm:"not null" json:"updateBy"`
	UpdateDt    time.Time                `gorm:"not null" json:"updateDt"`
	DeleteDt    gorm.DeletedAt           `gorm:"null" json:"deleteDt"`
}
