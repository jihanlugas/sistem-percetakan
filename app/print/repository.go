package print

import (
	"fmt"
	"github.com/jihanlugas/sistem-percetakan/model"
	"github.com/jihanlugas/sistem-percetakan/request"
	"github.com/jihanlugas/sistem-percetakan/utils"
	"gorm.io/gorm"
	"strings"
)

type Repository interface {
	GetTableById(conn *gorm.DB, id string, preloads ...string) (tPrint model.Print, err error)
	GetViewById(conn *gorm.DB, id string, preloads ...string) (vPrint model.PrintView, err error)
	Create(conn *gorm.DB, tPrint model.Print) error
	Creates(conn *gorm.DB, tPrints []model.Print) error
	Update(conn *gorm.DB, tPrint model.Print) error
	Save(conn *gorm.DB, tPrint model.Print) error
	Delete(conn *gorm.DB, tPrint model.Print) error
	DeleteByOrderId(conn *gorm.DB, id string) error
	Page(conn *gorm.DB, req request.PagePrint) (vPrints []model.PrintView, count int64, err error)
}

type repositoryy struct {
}

func (r repositoryy) GetTableById(conn *gorm.DB, id string, preloads ...string) (tPrint model.Print, err error) {
	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}
	err = conn.Where("id = ? ", id).First(&tPrint).Error
	return tPrint, err
}

func (r repositoryy) GetViewById(conn *gorm.DB, id string, preloads ...string) (vPrint model.PrintView, err error) {
	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}
	err = conn.Where("id = ? ", id).First(&vPrint).Error
	return vPrint, err
}

func (r repositoryy) Create(conn *gorm.DB, tPrint model.Print) error {
	return conn.Create(&tPrint).Error
}

func (r repositoryy) Creates(conn *gorm.DB, tPrints []model.Print) error {
	return conn.Create(&tPrints).Error
}

func (r repositoryy) Update(conn *gorm.DB, tPrint model.Print) error {
	return conn.Model(&tPrint).Updates(&tPrint).Error
}

func (r repositoryy) Save(conn *gorm.DB, tPrint model.Print) error {
	return conn.Save(&tPrint).Error
}

func (r repositoryy) Delete(conn *gorm.DB, tPrint model.Print) error {
	return conn.Delete(&tPrint).Error
}

func (r repositoryy) DeleteByOrderId(conn *gorm.DB, id string) error {
	return conn.Where("order_id = ? ", id).Delete(&model.Print{}).Error
}

func (r repositoryy) Page(conn *gorm.DB, req request.PagePrint) (vPrints []model.PrintView, count int64, err error) {
	query := conn.Model(&vPrints)

	// preloads
	preloads := strings.Split(req.Preloads, ",")
	for _, preload := range preloads {
		if utils.IsAvailablePreload(preload, model.PreloadPrint) {
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
	if req.StartDt != nil {
		query = query.Where("create_dt >= ?", req.StartDt)
	}
	if req.EndDt != nil {
		query = query.Where("create_dt <= ?", req.EndDt)
	}
	if req.StartTotalPrint != nil {
		query = query.Where("total >= ?", req.StartTotalPrint)
	}
	if req.EndTotalPrint != nil {
		query = query.Where("total <= ?", req.EndTotalPrint)
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

	if req.Limit >= 0 {
		query = query.Offset((req.GetPage() - 1) * req.GetLimit()).Limit(req.GetLimit())
	}

	err = query.Find(&vPrints).Error
	if err != nil {
		return vPrints, count, err
	}

	return vPrints, count, err
}

func NewRepository() Repository {
	return &repositoryy{}
}
