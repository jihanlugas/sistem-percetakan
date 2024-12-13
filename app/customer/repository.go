package customer

import (
	"fmt"
	"github.com/jihanlugas/sistem-percetakan/model"
	"github.com/jihanlugas/sistem-percetakan/request"
	"gorm.io/gorm"
)

type Repository interface {
	GetTableById(conn *gorm.DB, id string, preloads ...string) (tCustomer model.Customer, err error)
	GetViewById(conn *gorm.DB, id string, preloads ...string) (vCustomer model.CustomerView, err error)
	Create(conn *gorm.DB, tCustomer model.Customer) error
	Creates(conn *gorm.DB, tCustomers []model.Customer) error
	Update(conn *gorm.DB, tCustomer model.Customer) error
	Save(conn *gorm.DB, tCustomer model.Customer) error
	Delete(conn *gorm.DB, tCustomer model.Customer) error
	Page(conn *gorm.DB, req request.PageCustomer) (vCustomers []model.CustomerView, count int64, err error)
}

type repository struct {
}

func (r repository) GetTableById(conn *gorm.DB, id string, preloads ...string) (tCustomer model.Customer, err error) {
	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}
	err = conn.Where("id = ? ", id).First(&tCustomer).Error
	return tCustomer, err
}

func (r repository) GetViewById(conn *gorm.DB, id string, preloads ...string) (vCustomer model.CustomerView, err error) {
	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}
	err = conn.Where("id = ? ", id).First(&vCustomer).Error
	return vCustomer, err
}

func (r repository) Create(conn *gorm.DB, tCustomer model.Customer) error {
	return conn.Create(&tCustomer).Error
}

func (r repository) Creates(conn *gorm.DB, tCustomers []model.Customer) error {
	return conn.Create(&tCustomers).Error
}

func (r repository) Update(conn *gorm.DB, tCustomer model.Customer) error {
	return conn.Model(&tCustomer).Updates(&tCustomer).Error
}

func (r repository) Save(conn *gorm.DB, tCustomer model.Customer) error {
	return conn.Save(&tCustomer).Error
}

func (r repository) Delete(conn *gorm.DB, tCustomer model.Customer) error {
	return conn.Delete(&tCustomer).Error
}

func (r repository) Page(conn *gorm.DB, req request.PageCustomer) (vCustomers []model.CustomerView, count int64, err error) {
	query := conn.Model(&vCustomers)

	// preloads
	if req.Company {
		query = query.Preload("Company")
	}

	// query
	if req.CompanyID != "" {
		query = query.Where("company_id = ?", req.CompanyID)
	}
	if req.Name != "" {
		query = query.Where("name ILIKE ?", "%"+req.Name+"%")
	}
	if req.Description != "" {
		query = query.Where("description ILIKE ?", "%"+req.Description+"%")
	}
	if req.CompanyName != "" {
		query = query.Where("company_name ILIKE ?", "%"+req.CompanyName+"%")
	}
	if req.CreateName != "" {
		query = query.Where("create_name ILIKE ?", "%"+req.CreateName+"%")
	}

	err = query.Count(&count).Error
	if err != nil {
		return vCustomers, count, err
	}

	if req.SortField != "" {
		query = query.Order(fmt.Sprintf("%s %s", req.SortField, req.SortOrder))
	} else {
		query = query.Order(fmt.Sprintf("%s %s", "name", "asc"))
	}
	err = query.Offset((req.GetPage() - 1) * req.GetLimit()).
		Limit(req.GetLimit()).
		Find(&vCustomers).Error
	if err != nil {
		return vCustomers, count, err
	}

	return vCustomers, count, err
}

func NewRepository() Repository {
	return &repository{}
}
