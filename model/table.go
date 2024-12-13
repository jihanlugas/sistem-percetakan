package model

import (
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
	NoHp              string         `gorm:"not null" json:"noHp"`
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
	NoHp        string         `gorm:"not null" json:"noHp"`
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
	IsDone      bool           `gorm:"not null" json:"isDone"`
	CreateBy    string         `gorm:"not null" json:"createBy"`
	CreateDt    time.Time      `gorm:"not null" json:"createDt"`
	UpdateBy    string         `gorm:"not null" json:"updateBy"`
	UpdateDt    time.Time      `gorm:"not null" json:"updateDt"`
	DeleteDt    gorm.DeletedAt `gorm:"null" json:"deleteDt"`
}

type Paper struct {
	ID          string         `gorm:"primaryKey" json:"id"`
	CompanyID   string         `gorm:"not null" json:"companyId"`
	Name        string         `gorm:"not null" json:"name"`
	Description string         `gorm:"not null" json:"description"`
	CreateBy    string         `gorm:"not null" json:"createBy"`
	CreateDt    time.Time      `gorm:"not null" json:"createDt"`
	UpdateBy    string         `gorm:"not null" json:"updateBy"`
	UpdateDt    time.Time      `gorm:"not null" json:"updateDt"`
	DeleteDt    gorm.DeletedAt `gorm:"null" json:"deleteDt"`
}

type Design struct {
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

type Print struct {
	ID          string         `gorm:"primaryKey" json:"id"`
	CompanyID   string         `gorm:"not null" json:"companyId"`
	OrderID     string         `gorm:"not null" json:"orderId"`
	PaperID     string         `gorm:"not null" json:"paperId"`
	Name        string         `gorm:"not null" json:"name"`
	Description string         `gorm:"not null" json:"description"`
	IsDuplex    bool           `gorm:"not null" json:"isDuplex"`
	PageCount   int64          `gorm:"not null" json:"pageCount"`
	PrintCount  int64          `gorm:"not null" json:"printCount"`
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

type Other struct {
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
	Name      string         `gorm:"not null" json:"name"`
	CreateBy  string         `gorm:"not null" json:"createBy"`
	CreateDt  time.Time      `gorm:"not null" json:"createDt"`
	UpdateBy  string         `gorm:"not null" json:"updateBy"`
	UpdateDt  time.Time      `gorm:"not null" json:"updateDt"`
	DeleteDt  gorm.DeletedAt `gorm:"null" json:"deleteDt"`
}

type Payment struct {
	ID            string         `gorm:"primaryKey" json:"id"`
	CompanyID     string         `gorm:"not null" json:"companyId"`
	OrderID       string         `gorm:"not null" json:"orderId"`
	Name          string         `gorm:"not null" json:"name"`
	Description   string         `gorm:"not null" json:"description"`
	IsDonePayment bool           `gorm:"not null" json:"isDonePayment"`
	Amount        int64          `gorm:"not null" json:"amount"`
	CreateBy      string         `gorm:"not null" json:"createBy"`
	CreateDt      time.Time      `gorm:"not null" json:"createDt"`
	UpdateBy      string         `gorm:"not null" json:"updateBy"`
	UpdateDt      time.Time      `gorm:"not null" json:"updateDt"`
	DeleteDt      gorm.DeletedAt `gorm:"null" json:"deleteDt"`
}
