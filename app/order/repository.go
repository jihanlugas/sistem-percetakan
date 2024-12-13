package order

import (
	"fmt"
	"github.com/jihanlugas/sistem-percetakan/model"
	"github.com/jihanlugas/sistem-percetakan/request"
	"gorm.io/gorm"
)

type Repository interface {
	GetTableById(conn *gorm.DB, id string, preloads ...string) (tOrder model.Order, err error)
	GetViewById(conn *gorm.DB, id string, preloads ...string) (vOrder model.OrderView, err error)
	Create(conn *gorm.DB, tOrder model.Order) error
	Creates(conn *gorm.DB, tOrders []model.Order) error
	Update(conn *gorm.DB, tOrder model.Order) error
	Save(conn *gorm.DB, tOrder model.Order) error
	Delete(conn *gorm.DB, tOrder model.Order) error
	Page(conn *gorm.DB, req request.PageOrder) (vOrders []model.OrderView, count int64, err error)
}

type repository struct {
}

func (r repository) GetTableById(conn *gorm.DB, id string, preloads ...string) (tOrder model.Order, err error) {
	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}
	err = conn.Where("id = ? ", id).First(&tOrder).Error
	return tOrder, err
}

func (r repository) GetViewById(conn *gorm.DB, id string, preloads ...string) (vOrder model.OrderView, err error) {
	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}
	err = conn.Where("id = ? ", id).First(&vOrder).Error
	return vOrder, err
}

func (r repository) Create(conn *gorm.DB, tOrder model.Order) error {
	return conn.Create(&tOrder).Error
}

func (r repository) Creates(conn *gorm.DB, tOrders []model.Order) error {
	return conn.Create(&tOrders).Error
}

func (r repository) Update(conn *gorm.DB, tOrder model.Order) error {
	return conn.Model(&tOrder).Updates(&tOrder).Error
}

func (r repository) Save(conn *gorm.DB, tOrder model.Order) error {
	return conn.Save(&tOrder).Error
}

func (r repository) Delete(conn *gorm.DB, tOrder model.Order) error {
	return conn.Delete(&tOrder).Error
}

func (r repository) Page(conn *gorm.DB, req request.PageOrder) (vOrders []model.OrderView, count int64, err error) {
	query := conn.Model(&vOrders)

	if req.Company {
		query = query.Preload("Company")
	}
	if req.Customer {
		query = query.Preload("Customer")
	}

	if req.CompanyID != "" {
		query = query.Where("company_id = ?", req.CompanyID)
	}
	if req.CustomerID != "" {
		query = query.Where("customer_id = ?", req.CustomerID)
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
	if req.CustomerName != "" {
		query = query.Where("customer_name ILIKE ?", "%"+req.CustomerName+"%")
	}
	if req.CreateName != "" {
		query = query.Where("create_name ILIKE ?", "%"+req.CreateName+"%")
	}

	err = query.Count(&count).Error
	if err != nil {
		return vOrders, count, err
	}

	if req.SortField != "" {
		query = query.Order(fmt.Sprintf("%s %s", req.SortField, req.SortOrder))
	} else {
		query = query.Order(fmt.Sprintf("%s %s", "name", "asc"))
	}
	err = query.Offset((req.GetPage() - 1) * req.GetLimit()).
		Limit(req.GetLimit()).
		Find(&vOrders).Error
	if err != nil {
		return vOrders, count, err
	}

	return vOrders, count, err
}

func NewRepository() Repository {
	return &repository{}
}
