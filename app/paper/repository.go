package paper

import (
	"fmt"
	"github.com/jihanlugas/sistem-percetakan/model"
	"github.com/jihanlugas/sistem-percetakan/request"
	"github.com/jihanlugas/sistem-percetakan/utils"
	"gorm.io/gorm"
	"strings"
)

type Repository interface {
	GetTableById(conn *gorm.DB, id string, preloads ...string) (tPaper model.Paper, err error)
	GetViewById(conn *gorm.DB, id string, preloads ...string) (vPaper model.PaperView, err error)
	Create(conn *gorm.DB, tPaper model.Paper) error
	Creates(conn *gorm.DB, tPapers []model.Paper) error
	Update(conn *gorm.DB, tPaper model.Paper) error
	Save(conn *gorm.DB, tPaper model.Paper) error
	Delete(conn *gorm.DB, tPaper model.Paper) error
	Page(conn *gorm.DB, req request.PagePaper) (vPapers []model.PaperView, count int64, err error)
}

type repository struct {
}

func (r repository) GetTableById(conn *gorm.DB, id string, preloads ...string) (tPaper model.Paper, err error) {
	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}
	err = conn.Where("id = ? ", id).First(&tPaper).Error
	return tPaper, err
}

func (r repository) GetViewById(conn *gorm.DB, id string, preloads ...string) (vPaper model.PaperView, err error) {
	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}
	err = conn.Where("id = ? ", id).First(&vPaper).Error
	return vPaper, err
}

func (r repository) Create(conn *gorm.DB, tPaper model.Paper) error {
	return conn.Create(&tPaper).Error
}

func (r repository) Creates(conn *gorm.DB, tPapers []model.Paper) error {
	return conn.Create(&tPapers).Error
}

func (r repository) Update(conn *gorm.DB, tPaper model.Paper) error {
	return conn.Model(&tPaper).Updates(&tPaper).Error
}

func (r repository) Save(conn *gorm.DB, tPaper model.Paper) error {
	return conn.Save(&tPaper).Error
}

func (r repository) Delete(conn *gorm.DB, tPaper model.Paper) error {
	return conn.Delete(&tPaper).Error
}

func (r repository) Page(conn *gorm.DB, req request.PagePaper) (vPapers []model.PaperView, count int64, err error) {
	query := conn.Model(&vPapers)

	// preloads
	preloads := strings.Split(req.Preloads, ",")
	for _, preload := range preloads {
		if utils.IsAvailablePreload(preload, model.PreloadPaper) {
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
		return vPapers, count, err
	}

	if req.SortField != "" {
		query = query.Order(fmt.Sprintf("%s %s", req.SortField, req.SortOrder))
	} else {
		query = query.Order(fmt.Sprintf("%s %s", "name", "asc"))
	}

	if req.Limit >= 0 {
		query = query.Offset((req.GetPage() - 1) * req.GetLimit()).Limit(req.GetLimit())
	}

	err = query.Find(&vPapers).Error
	if err != nil {
		return vPapers, count, err
	}

	return vPapers, count, err
}

func NewRepository() Repository {
	return &repository{}
}
