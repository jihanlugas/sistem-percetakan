package order

import (
	"fmt"
	"github.com/jihanlugas/sistem-percetakan/model"
	"github.com/jihanlugas/sistem-percetakan/request"
	"github.com/jihanlugas/sistem-percetakan/utils"
	"gorm.io/gorm"
	"strings"
)

type Repository interface {
	GetTableById(conn *gorm.DB, id string, preloads ...string) (tOrder model.Order, err error)
	GetViewById(conn *gorm.DB, id string, preloads ...string) (vOrder model.OrderView, err error)
	GetNextNumber(conn *gorm.DB, companyID string) (number int64)
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
		if utils.IsAvailablePreload(preload, model.PreloadOrder) {
			conn = conn.Preload(preload)
		}
	}
	err = conn.Where("id = ? ", id).First(&tOrder).Error
	return tOrder, err
}

func (r repository) GetViewById(conn *gorm.DB, id string, preloads ...string) (vOrder model.OrderView, err error) {
	for _, preload := range preloads {
		if utils.IsAvailablePreload(preload, model.PreloadOrder) {
			conn = conn.Preload(preload)
		}
	}
	err = conn.Where("id = ? ", id).First(&vOrder).Error
	return vOrder, err
}

func (r repository) GetNextNumber(conn *gorm.DB, companyID string) (number int64) {
	conn.Model(&model.Order{}).Where("company_id = ?", companyID).Count(&number)
	return number + 1
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

	preloads := strings.Split(req.Preloads, ",")
	for _, preload := range preloads {
		if utils.IsAvailablePreload(preload, model.PreloadOrder) {
			query = query.Preload(preload)
		}
	}

	if req.CompanyID != "" {
		query = query.Where("company_id = ?", req.CompanyID)
	}
	if req.CustomerID != "" {
		query = query.Where("customer_id = ?", req.CustomerID)
	}
	if req.PhaseID != "" {
		query = query.Where("phase_id = ?", req.PhaseID)
	}
	if req.Name != "" {
		query = query.Where("name ILIKE ?", "%"+req.Name+"%")
	}
	if req.Description != "" {
		query = query.Where("description ILIKE ?", "%"+req.Description+"%")
	}
	if req.StartDt != nil {
		query = query.Where("create_dt >= ?", req.StartDt)
	}
	if req.EndDt != nil {
		query = query.Where("create_dt <= ?", req.EndDt)
	}
	if req.StartTotalOrder != nil {
		query = query.Where("total_order >= ?", req.StartTotalOrder)
	}
	if req.EndTotalOrder != nil {
		query = query.Where("total_order <= ?", req.EndTotalOrder)
	}

	err = query.Count(&count).Error
	if err != nil {
		return vOrders, count, err
	}

	if req.SortField != "" {
		query = query.Order(fmt.Sprintf("%s %s", req.SortField, req.SortOrder))
	} else {
		query = query.Order(fmt.Sprintf("%s %s", "create_dt", "desc"))
	}

	if req.Limit >= 0 {
		query = query.Offset((req.GetPage() - 1) * req.GetLimit()).Limit(req.GetLimit())
	}

	err = query.Find(&vOrders).Error
	if err != nil {
		return vOrders, count, err
	}

	return vOrders, count, err
}

func NewRepository() Repository {
	return &repository{}
}
