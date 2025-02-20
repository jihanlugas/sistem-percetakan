package transaction

import (
	"fmt"
	"github.com/jihanlugas/sistem-percetakan/model"
	"github.com/jihanlugas/sistem-percetakan/request"
	"github.com/jihanlugas/sistem-percetakan/utils"
	"gorm.io/gorm"
	"strings"
)

type Repository interface {
	GetTableById(conn *gorm.DB, id string, preloads ...string) (tTransaction model.Transaction, err error)
	GetViewById(conn *gorm.DB, id string, preloads ...string) (vTransaction model.TransactionView, err error)
	Create(conn *gorm.DB, tTransaction model.Transaction) error
	Creates(conn *gorm.DB, tTransactions []model.Transaction) error
	Update(conn *gorm.DB, tTransaction model.Transaction) error
	Save(conn *gorm.DB, tTransaction model.Transaction) error
	Delete(conn *gorm.DB, tTransaction model.Transaction) error
	DeleteByOrderId(conn *gorm.DB, id string) error
	Page(conn *gorm.DB, req request.PageTransaction) (vTransactions []model.TransactionView, count int64, err error)
}

type repository struct {
}

func (r repository) GetTableById(conn *gorm.DB, id string, preloads ...string) (tTransaction model.Transaction, err error) {
	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}
	err = conn.Where("id = ? ", id).First(&tTransaction).Error
	return tTransaction, err
}

func (r repository) GetViewById(conn *gorm.DB, id string, preloads ...string) (vTransaction model.TransactionView, err error) {
	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}
	err = conn.Where("id = ? ", id).First(&vTransaction).Error
	return vTransaction, err
}

func (r repository) Create(conn *gorm.DB, tTransaction model.Transaction) error {
	return conn.Create(&tTransaction).Error
}

func (r repository) Creates(conn *gorm.DB, tTransactions []model.Transaction) error {
	return conn.Create(&tTransactions).Error
}

func (r repository) Update(conn *gorm.DB, tTransaction model.Transaction) error {
	return conn.Model(&tTransaction).Updates(&tTransaction).Error
}

func (r repository) Save(conn *gorm.DB, tTransaction model.Transaction) error {
	return conn.Save(&tTransaction).Error
}

func (r repository) Delete(conn *gorm.DB, tTransaction model.Transaction) error {
	return conn.Delete(&tTransaction).Error
}

func (r repository) DeleteByOrderId(conn *gorm.DB, id string) error {
	return conn.Where("order_id = ? ", id).Delete(&model.Transaction{}).Error
}

func (r repository) Page(conn *gorm.DB, req request.PageTransaction) (vTransactions []model.TransactionView, count int64, err error) {
	query := conn.Model(&vTransactions)

	// preloads
	preloads := strings.Split(req.Preloads, ",")
	for _, preload := range preloads {
		if utils.IsAvailablePreload(preload, model.PreloadTransaction) {
			query = query.Preload(preload)
		}
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
	if req.Type != 0 {
		query = query.Where("type = ?", req.Type)
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
	if req.StartDt != nil {
		query = query.Where("create_dt >= ?", req.StartDt)
	}
	if req.EndDt != nil {
		query = query.Where("create_dt <= ?", req.EndDt)
	}
	if req.StartAmount != nil {
		query = query.Where("amount >= ?", req.StartAmount)
	}
	if req.EndAmount != nil {
		query = query.Where("amount <= ?", req.EndAmount)
	}

	err = query.Count(&count).Error
	if err != nil {
		return vTransactions, count, err
	}

	if req.SortField != "" {
		query = query.Order(fmt.Sprintf("%s %s", req.SortField, req.SortOrder))
	} else {
		query = query.Order(fmt.Sprintf("%s %s", "name", "asc"))
	}

	if req.Limit >= 0 {
		query = query.Offset((req.GetPage() - 1) * req.GetLimit()).Limit(req.GetLimit())
	}

	err = query.Find(&vTransactions).Error
	if err != nil {
		return vTransactions, count, err
	}

	return vTransactions, count, err
}

func NewRepository() Repository {
	return &repository{}
}
