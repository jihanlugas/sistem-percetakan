package print

import (
	"fmt"
	"github.com/jihanlugas/sistem-percetakan/model"
	"github.com/jihanlugas/sistem-percetakan/request"
	"gorm.io/gorm"
)

type Repository interface {
	GetTableById(conn *gorm.DB, id string, preloads ...string) (tPrint model.Print, err error)
	GetViewById(conn *gorm.DB, id string, preloads ...string) (vPrint model.PrintView, err error)
	Create(conn *gorm.DB, tPrint model.Print) error
	Creates(conn *gorm.DB, tPrints []model.Print) error
	Update(conn *gorm.DB, tPrint model.Print) error
	Save(conn *gorm.DB, tPrint model.Print) error
	Delete(conn *gorm.DB, tPrint model.Print) error
	Page(conn *gorm.DB, req request.PagePrint) (vPrints []model.PrintView, count int64, err error)
}

type repository struct {
}

func (r repository) GetTableById(conn *gorm.DB, id string, preloads ...string) (tPrint model.Print, err error) {
	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}
	err = conn.Where("id = ? ", id).First(&tPrint).Error
	return tPrint, err
}

func (r repository) GetViewById(conn *gorm.DB, id string, preloads ...string) (vPrint model.PrintView, err error) {
	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}
	err = conn.Where("id = ? ", id).First(&vPrint).Error
	return vPrint, err
}

func (r repository) Create(conn *gorm.DB, tPrint model.Print) error {
	return conn.Create(&tPrint).Error
}

func (r repository) Creates(conn *gorm.DB, tPrints []model.Print) error {
	return conn.Create(&tPrints).Error
}

func (r repository) Update(conn *gorm.DB, tPrint model.Print) error {
	return conn.Model(&tPrint).Updates(&tPrint).Error
}

func (r repository) Save(conn *gorm.DB, tPrint model.Print) error {
	return conn.Save(&tPrint).Error
}

func (r repository) Delete(conn *gorm.DB, tPrint model.Print) error {
	return conn.Delete(&tPrint).Error
}

func (r repository) Page(conn *gorm.DB, req request.PagePrint) (vPrints []model.PrintView, count int64, err error) {
	query := conn.Model(&vPrints)

	// preloads
	if req.Company {
		query = query.Preload("Company")
	}
	if req.Order {
		query = query.Preload("Order")
	}
	if req.Paper {
		query = query.Preload("Paper")
	}

	// query
	if req.CompanyID != "" {
		query = query.Where("company_id = ?", req.CompanyID)
	}
	if req.OrderID != "" {
		query = query.Where("order_id = ?", req.OrderID)
	}
	if req.PaperID != "" {
		query = query.Where("paper_id = ?", req.PaperID)
	}
	if req.Name != "" {
		query = query.Where("name ILIKE ?", "%"+req.Name+"%")
	}
	if req.Description != "" {
		query = query.Where("description ILIKE ?", "%"+req.Description+"%")
	}
	if req.IsDuplex != nil {
		query = query.Where("is_duplex = ?", req.Description)
	}
	if req.CompanyName != "" {
		query = query.Where("company_name ILIKE ?", "%"+req.CompanyName+"%")
	}
	if req.OrderName != "" {
		query = query.Where("order_name ILIKE ?", "%"+req.OrderName+"%")
	}
	if req.PaperName != "" {
		query = query.Where("paper_name ILIKE ?", "%"+req.PaperName+"%")
	}
	if req.CreateName != "" {
		query = query.Where("create_name ILIKE ?", "%"+req.CreateName+"%")
	}

	err = query.Count(&count).Error
	if err != nil {
		return vPrints, count, err
	}

	if req.SortField != "" {
		query = query.Order(fmt.Sprintf("%s %s", req.SortField, req.SortOrder))
	} else {
		query = query.Order(fmt.Sprintf("%s %s", "name", "asc"))
	}
	err = query.Offset((req.GetPage() - 1) * req.GetLimit()).
		Limit(req.GetLimit()).
		Find(&vPrints).Error
	if err != nil {
		return vPrints, count, err
	}

	return vPrints, count, err
}

func NewRepository() Repository {
	return &repository{}
}
