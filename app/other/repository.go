package other

import (
	"fmt"
	"github.com/jihanlugas/sistem-percetakan/model"
	"github.com/jihanlugas/sistem-percetakan/request"
	"gorm.io/gorm"
)

type Repository interface {
	GetTableById(conn *gorm.DB, id string, preloads ...string) (tOther model.Other, err error)
	GetViewById(conn *gorm.DB, id string, preloads ...string) (vOther model.OtherView, err error)
	Create(conn *gorm.DB, tOther model.Other) error
	Creates(conn *gorm.DB, tOthers []model.Other) error
	Update(conn *gorm.DB, tOther model.Other) error
	Save(conn *gorm.DB, tOther model.Other) error
	Delete(conn *gorm.DB, tOther model.Other) error
	Page(conn *gorm.DB, req request.PageOther) (vOthers []model.OtherView, count int64, err error)
}

type repository struct {
}

func (r repository) GetTableById(conn *gorm.DB, id string, preloads ...string) (tOther model.Other, err error) {
	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}
	err = conn.Where("id = ? ", id).First(&tOther).Error
	return tOther, err
}

func (r repository) GetViewById(conn *gorm.DB, id string, preloads ...string) (vOther model.OtherView, err error) {
	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}
	err = conn.Where("id = ? ", id).First(&vOther).Error
	return vOther, err
}

func (r repository) Create(conn *gorm.DB, tOther model.Other) error {
	return conn.Create(&tOther).Error
}

func (r repository) Creates(conn *gorm.DB, tOthers []model.Other) error {
	return conn.Create(&tOthers).Error
}

func (r repository) Update(conn *gorm.DB, tOther model.Other) error {
	return conn.Model(&tOther).Updates(&tOther).Error
}

func (r repository) Save(conn *gorm.DB, tOther model.Other) error {
	return conn.Save(&tOther).Error
}

func (r repository) Delete(conn *gorm.DB, tOther model.Other) error {
	return conn.Delete(&tOther).Error
}

func (r repository) Page(conn *gorm.DB, req request.PageOther) (vOthers []model.OtherView, count int64, err error) {
	query := conn.Model(&vOthers)

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
		return vOthers, count, err
	}

	if req.SortField != "" {
		query = query.Order(fmt.Sprintf("%s %s", req.SortField, req.SortOrder))
	} else {
		query = query.Order(fmt.Sprintf("%s %s", "name", "asc"))
	}
	err = query.Offset((req.GetPage() - 1) * req.GetLimit()).
		Limit(req.GetLimit()).
		Find(&vOthers).Error
	if err != nil {
		return vOthers, count, err
	}

	return vOthers, count, err
}

func NewRepository() Repository {
	return &repository{}
}
