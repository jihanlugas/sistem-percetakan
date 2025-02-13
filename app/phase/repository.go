package phase

import (
	"fmt"
	"github.com/jihanlugas/sistem-percetakan/model"
	"github.com/jihanlugas/sistem-percetakan/request"
	"github.com/jihanlugas/sistem-percetakan/utils"
	"gorm.io/gorm"
	"strings"
)

type Repository interface {
	GetTableById(conn *gorm.DB, id string, preloads ...string) (tPhase model.Phase, err error)
	GetViewById(conn *gorm.DB, id string, preloads ...string) (vPhase model.PhaseView, err error)
	Create(conn *gorm.DB, tPhase model.Phase) error
	Creates(conn *gorm.DB, tPhases []model.Phase) error
	Update(conn *gorm.DB, tPhase model.Phase) error
	Save(conn *gorm.DB, tPhase model.Phase) error
	Delete(conn *gorm.DB, tPhase model.Phase) error
	Page(conn *gorm.DB, req request.PagePhase) (vPhases []model.PhaseView, count int64, err error)
}

type repository struct {
}

func (r repository) GetTableById(conn *gorm.DB, id string, preloads ...string) (tPhase model.Phase, err error) {
	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}
	err = conn.Where("id = ? ", id).First(&tPhase).Error
	return tPhase, err
}

func (r repository) GetViewById(conn *gorm.DB, id string, preloads ...string) (vPhase model.PhaseView, err error) {
	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}
	err = conn.Where("id = ? ", id).First(&vPhase).Error
	return vPhase, err
}

func (r repository) Create(conn *gorm.DB, tPhase model.Phase) error {
	return conn.Create(&tPhase).Error
}

func (r repository) Creates(conn *gorm.DB, tPhases []model.Phase) error {
	return conn.Create(&tPhases).Error
}

func (r repository) Update(conn *gorm.DB, tPhase model.Phase) error {
	return conn.Model(&tPhase).Updates(&tPhase).Error
}

func (r repository) Save(conn *gorm.DB, tPhase model.Phase) error {
	return conn.Save(&tPhase).Error
}

func (r repository) Delete(conn *gorm.DB, tPhase model.Phase) error {
	return conn.Delete(&tPhase).Error
}

func (r repository) Page(conn *gorm.DB, req request.PagePhase) (vPhases []model.PhaseView, count int64, err error) {
	query := conn.Model(&vPhases)

	// preloads
	preloads := strings.Split(req.Preloads, ",")
	for _, preload := range preloads {
		if utils.IsAvailablePreload(preload, model.PreloadPhase) {
			query = query.Preload(preload)
		}
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
		return vPhases, count, err
	}

	if req.SortField != "" {
		query = query.Order(fmt.Sprintf("%s %s", req.SortField, req.SortOrder))
	} else {
		query = query.Order(fmt.Sprintf("%s %s", "\"order\"", "asc"))
	}

	if req.Limit >= 0 {
		query = query.Offset((req.GetPage() - 1) * req.GetLimit()).Limit(req.GetLimit())
	}

	err = query.Find(&vPhases).Error
	if err != nil {
		return vPhases, count, err
	}

	return vPhases, count, err
}

func NewRepository() Repository {
	return &repository{}
}
