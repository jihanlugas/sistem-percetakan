package payment

import (
	"fmt"
	"github.com/jihanlugas/sistem-percetakan/model"
	"github.com/jihanlugas/sistem-percetakan/request"
	"gorm.io/gorm"
)

type Repository interface {
	GetTableById(conn *gorm.DB, id string, preloads ...string) (tPayment model.Payment, err error)
	GetViewById(conn *gorm.DB, id string, preloads ...string) (vPayment model.PaymentView, err error)
	Create(conn *gorm.DB, tPayment model.Payment) error
	Creates(conn *gorm.DB, tPayments []model.Payment) error
	Update(conn *gorm.DB, tPayment model.Payment) error
	Save(conn *gorm.DB, tPayment model.Payment) error
	Delete(conn *gorm.DB, tPayment model.Payment) error
	Page(conn *gorm.DB, req request.PagePayment) (vPayments []model.PaymentView, count int64, err error)
}

type repository struct {
}

func (r repository) GetTableById(conn *gorm.DB, id string, preloads ...string) (tPayment model.Payment, err error) {
	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}
	err = conn.Where("id = ? ", id).First(&tPayment).Error
	return tPayment, err
}

func (r repository) GetViewById(conn *gorm.DB, id string, preloads ...string) (vPayment model.PaymentView, err error) {
	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}
	err = conn.Where("id = ? ", id).First(&vPayment).Error
	return vPayment, err
}

func (r repository) Create(conn *gorm.DB, tPayment model.Payment) error {
	return conn.Create(&tPayment).Error
}

func (r repository) Creates(conn *gorm.DB, tPayments []model.Payment) error {
	return conn.Create(&tPayments).Error
}

func (r repository) Update(conn *gorm.DB, tPayment model.Payment) error {
	return conn.Model(&tPayment).Updates(&tPayment).Error
}

func (r repository) Save(conn *gorm.DB, tPayment model.Payment) error {
	return conn.Save(&tPayment).Error
}

func (r repository) Delete(conn *gorm.DB, tPayment model.Payment) error {
	return conn.Delete(&tPayment).Error
}

func (r repository) Page(conn *gorm.DB, req request.PagePayment) (vPayments []model.PaymentView, count int64, err error) {
	query := conn.Model(&vPayments)

	// preloads
	if req.Company {
		query = query.Preload("Company")
	}
	if req.Order {
		query = query.Preload("Order")
	}

	// query
	if req.CompanyID != "" {
		query = query.Where("company_id = ?", req.CompanyID)
	}
	if req.OrderID != "" {
		query = query.Where("order_id = ?", req.OrderID)
	}
	if req.Name != "" {
		query = query.Where("name ILIKE ?", "%"+req.Name+"%")
	}
	if req.Description != "" {
		query = query.Where("description ILIKE ?", "%"+req.Description+"%")
	}
	if req.IsDonePayment != nil {
		query = query.Where("is_done_payment = ?", req.IsDonePayment)
	}
	if req.CompanyName != "" {
		query = query.Where("company_name ILIKE ?", "%"+req.CompanyName+"%")
	}
	if req.OrderName != "" {
		query = query.Where("order_name ILIKE ?", "%"+req.OrderName+"%")
	}
	if req.CreateName != "" {
		query = query.Where("create_name ILIKE ?", "%"+req.CreateName+"%")
	}

	err = query.Count(&count).Error
	if err != nil {
		return vPayments, count, err
	}

	if req.SortField != "" {
		query = query.Order(fmt.Sprintf("%s %s", req.SortField, req.SortOrder))
	} else {
		query = query.Order(fmt.Sprintf("%s %s", "name", "asc"))
	}
	err = query.Offset((req.GetPage() - 1) * req.GetLimit()).
		Limit(req.GetLimit()).
		Find(&vPayments).Error
	if err != nil {
		return vPayments, count, err
	}

	return vPayments, count, err
}

func NewRepository() Repository {
	return &repository{}
}
