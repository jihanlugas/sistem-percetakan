package company

import (
	"fmt"
	"github.com/jihanlugas/sistem-percetakan/model"
	"github.com/jihanlugas/sistem-percetakan/request"
	"gorm.io/gorm"
)

type Repository interface {
	GetTableById(conn *gorm.DB, id string) (tCompany model.Company, err error)
	GetTableByName(conn *gorm.DB, name string) (tCompany model.Company, err error)
	GetViewById(conn *gorm.DB, id string) (vCompany model.CompanyView, err error)
	GetViewByName(conn *gorm.DB, name string) (vCompany model.CompanyView, err error)
	Create(conn *gorm.DB, tCompany model.Company) error
	Update(conn *gorm.DB, tCompany model.Company) error
	Save(conn *gorm.DB, tCompany model.Company) error
	Delete(conn *gorm.DB, tCompany model.Company) error
	Page(conn *gorm.DB, req request.PageCompany) (vCompanies []model.CompanyView, count int64, err error)
}

type repository struct {
}

func (r repository) GetTableById(conn *gorm.DB, id string) (tCompany model.Company, err error) {
	err = conn.Where("id = ? ", id).First(&tCompany).Error
	return tCompany, err
}

func (r repository) GetTableByName(conn *gorm.DB, name string) (tCompany model.Company, err error) {
	err = conn.Where("name = ? ", name).First(&tCompany).Error
	return tCompany, err
}

func (r repository) GetViewById(conn *gorm.DB, id string) (vCompany model.CompanyView, err error) {
	err = conn.Where("id = ? ", id).First(&vCompany).Error
	return vCompany, err
}

func (r repository) GetViewByName(conn *gorm.DB, name string) (vCompany model.CompanyView, err error) {
	err = conn.Where("name = ? ", name).First(&vCompany).Error
	return vCompany, err
}

func (r repository) Create(conn *gorm.DB, tCompany model.Company) error {
	return conn.Create(&tCompany).Error
}

func (r repository) Update(conn *gorm.DB, tCompany model.Company) error {
	return conn.Model(&tCompany).Updates(&tCompany).Error
}

func (r repository) Save(conn *gorm.DB, tCompany model.Company) error {
	return conn.Save(&tCompany).Error
}

func (r repository) Delete(conn *gorm.DB, tCompany model.Company) error {
	return conn.Delete(&tCompany).Error
}

func (r repository) Page(conn *gorm.DB, req request.PageCompany) (vCompanies []model.CompanyView, count int64, err error) {
	query := conn.Model(&vCompanies)
	if req.Name != "" {
		query = query.Where("name ILIKE ?", "%"+req.Name+"%")
	}
	if req.Description != "" {
		query = query.Where("description ILIKE ?", "%"+req.Description+"%")
	}
	if req.CreateName != "" {
		query = query.Where("create_name ILIKE ?", "%"+req.CreateName+"%")
	}

	err = query.Count(&count).Error
	if err != nil {
		return vCompanies, count, err
	}

	if req.SortField != "" {
		query = query.Order(fmt.Sprintf("%s %s", req.SortField, req.SortOrder))
	} else {
		query = query.Order(fmt.Sprintf("%s %s", "name", "asc"))
	}
	if req.Limit >= 0 {
		query = query.Offset((req.GetPage() - 1) * req.GetLimit()).Limit(req.GetLimit())
	}

	err = query.Find(&vCompanies).Error
	if err != nil {
		return vCompanies, count, err
	}

	return vCompanies, count, err
}

func NewRepository() Repository {
	return repository{}
}
