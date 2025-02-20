package finishing

import (
	"fmt"
	"github.com/jihanlugas/sistem-percetakan/model"
	"github.com/jihanlugas/sistem-percetakan/request"
	"github.com/jihanlugas/sistem-percetakan/utils"
	"gorm.io/gorm"
	"strings"
)

type Repository interface {
	GetTableById(conn *gorm.DB, id string, preloads ...string) (tFinishing model.Finishing, err error)
	GetViewById(conn *gorm.DB, id string, preloads ...string) (vFinishing model.FinishingView, err error)
	Create(conn *gorm.DB, tFinishing model.Finishing) error
	Creates(conn *gorm.DB, tFinishings []model.Finishing) error
	Update(conn *gorm.DB, tFinishing model.Finishing) error
	Save(conn *gorm.DB, tFinishing model.Finishing) error
	Delete(conn *gorm.DB, tFinishing model.Finishing) error
	DeleteByOrderId(conn *gorm.DB, id string) error
	Page(conn *gorm.DB, req request.PageFinishing) (vFinishings []model.FinishingView, count int64, err error)
}

type repository struct {
}

func (r repository) GetTableById(conn *gorm.DB, id string, preloads ...string) (tFinishing model.Finishing, err error) {
	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}
	err = conn.Where("id = ? ", id).First(&tFinishing).Error
	return tFinishing, err
}

func (r repository) GetViewById(conn *gorm.DB, id string, preloads ...string) (vFinishing model.FinishingView, err error) {
	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}
	err = conn.Where("id = ? ", id).First(&vFinishing).Error
	return vFinishing, err
}

func (r repository) Create(conn *gorm.DB, tFinishing model.Finishing) error {
	return conn.Create(&tFinishing).Error
}

func (r repository) Creates(conn *gorm.DB, tFinishings []model.Finishing) error {
	return conn.Create(&tFinishings).Error
}

func (r repository) Update(conn *gorm.DB, tFinishing model.Finishing) error {
	return conn.Model(&tFinishing).Updates(&tFinishing).Error
}

func (r repository) Save(conn *gorm.DB, tFinishing model.Finishing) error {
	return conn.Save(&tFinishing).Error
}

func (r repository) Delete(conn *gorm.DB, tFinishing model.Finishing) error {
	return conn.Delete(&tFinishing).Error
}

func (r repository) DeleteByOrderId(conn *gorm.DB, id string) error {
	return conn.Where("order_id = ? ", id).Delete(&model.Finishing{}).Error
}

func (r repository) Page(conn *gorm.DB, req request.PageFinishing) (vFinishings []model.FinishingView, count int64, err error) {
	query := conn.Model(&vFinishings)

	// preloads
	preloads := strings.Split(req.Preloads, ",")
	for _, preload := range preloads {
		if utils.IsAvailablePreload(preload, model.PreloadFinishing) {
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
	if req.StartTotalFinishing != nil {
		query = query.Where("total >= ?", req.StartTotalFinishing)
	}
	if req.EndTotalFinishing != nil {
		query = query.Where("total <= ?", req.EndTotalFinishing)
	}

	err = query.Count(&count).Error
	if err != nil {
		return vFinishings, count, err
	}

	if req.SortField != "" {
		query = query.Order(fmt.Sprintf("%s %s", req.SortField, req.SortOrder))
	} else {
		query = query.Order(fmt.Sprintf("%s %s", "name", "asc"))
	}
	if req.Limit >= 0 {
		query = query.Offset((req.GetPage() - 1) * req.GetLimit()).Limit(req.GetLimit())
	}

	err = query.Find(&vFinishings).Error
	if err != nil {
		return vFinishings, count, err
	}

	return vFinishings, count, err
}

func NewRepository() Repository {
	return &repository{}
}
