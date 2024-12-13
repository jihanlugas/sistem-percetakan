package design

import (
	"fmt"
	"github.com/jihanlugas/sistem-percetakan/model"
	"github.com/jihanlugas/sistem-percetakan/request"
	"gorm.io/gorm"
)

type Repository interface {
	GetTableById(conn *gorm.DB, id string, preloads ...string) (tDesign model.Design, err error)
	GetViewById(conn *gorm.DB, id string, preloads ...string) (vDesign model.DesignView, err error)
	Create(conn *gorm.DB, tDesign model.Design) error
	Creates(conn *gorm.DB, tDesigns []model.Design) error
	Update(conn *gorm.DB, tDesign model.Design) error
	Save(conn *gorm.DB, tDesign model.Design) error
	Delete(conn *gorm.DB, tDesign model.Design) error
	Page(conn *gorm.DB, req request.PageDesign) (vDesigns []model.DesignView, count int64, err error)
}

type repository struct {
}

func (r repository) GetTableById(conn *gorm.DB, id string, preloads ...string) (tDesign model.Design, err error) {
	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}
	err = conn.Where("id = ? ", id).First(&tDesign).Error
	return tDesign, err
}

func (r repository) GetViewById(conn *gorm.DB, id string, preloads ...string) (vDesign model.DesignView, err error) {
	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}
	err = conn.Where("id = ? ", id).First(&vDesign).Error
	return vDesign, err
}

func (r repository) Create(conn *gorm.DB, tDesign model.Design) error {
	return conn.Create(&tDesign).Error
}

func (r repository) Creates(conn *gorm.DB, tDesigns []model.Design) error {
	return conn.Create(&tDesigns).Error
}

func (r repository) Update(conn *gorm.DB, tDesign model.Design) error {
	return conn.Model(&tDesign).Updates(&tDesign).Error
}

func (r repository) Save(conn *gorm.DB, tDesign model.Design) error {
	return conn.Save(&tDesign).Error
}

func (r repository) Delete(conn *gorm.DB, tDesign model.Design) error {
	return conn.Delete(&tDesign).Error
}

func (r repository) Page(conn *gorm.DB, req request.PageDesign) (vDesigns []model.DesignView, count int64, err error) {
	query := conn.Model(&vDesigns)

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
		return vDesigns, count, err
	}

	if req.SortField != "" {
		query = query.Order(fmt.Sprintf("%s %s", req.SortField, req.SortOrder))
	} else {
		query = query.Order(fmt.Sprintf("%s %s", "name", "asc"))
	}
	err = query.Offset((req.GetPage() - 1) * req.GetLimit()).
		Limit(req.GetLimit()).
		Find(&vDesigns).Error
	if err != nil {
		return vDesigns, count, err
	}

	return vDesigns, count, err
}

func NewRepository() Repository {
	return &repository{}
}
