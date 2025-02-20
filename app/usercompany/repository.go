package usercompany

import (
	"fmt"
	"github.com/jihanlugas/sistem-percetakan/model"
	"github.com/jihanlugas/sistem-percetakan/request"
	"gorm.io/gorm"
)

type Repository interface {
	GetTableById(conn *gorm.DB, id string) (tUsercompany model.Usercompany, err error)
	GetCreatorByCompanyId(conn *gorm.DB, companyID string) (tUsercompany model.Usercompany, err error)
	GetCompanyDefaultByUserId(conn *gorm.DB, userID string) (tUsercompany model.Usercompany, err error)
	GetViewById(conn *gorm.DB, id string) (vUsercompany model.UsercompanyView, err error)
	GetViewByUserIdAndCompanyId(conn *gorm.DB, userID, companyID string) (vUsercompany model.UsercompanyView, err error)
	GetViewCreatorByCompanyId(conn *gorm.DB, companyID string) (vUsercompany model.UsercompanyView, err error)
	GetViewCompanyDefaultByUserId(conn *gorm.DB, userID string) (vUsercompany model.UsercompanyView, err error)
	Create(conn *gorm.DB, tUsercompany model.Usercompany) error
	Update(conn *gorm.DB, tUsercompany model.Usercompany) error
	Save(conn *gorm.DB, tUsercompany model.Usercompany) error
	Delete(conn *gorm.DB, tUsercompany model.Usercompany) error
	Page(conn *gorm.DB, req request.PageUsercompany) (tUsercompanies []model.UsercompanyView, count int64, err error)
}

type repository struct {
}

func (r repository) GetTableById(conn *gorm.DB, id string) (tUsercompany model.Usercompany, err error) {
	err = conn.Where("id = ? ", id).First(&tUsercompany).Error
	return tUsercompany, err
}

func (r repository) GetCreatorByCompanyId(conn *gorm.DB, companyID string) (tUsercompany model.Usercompany, err error) {
	err = conn.Where("company_id = ? ", companyID).
		Where("is_creator = ? ", true).
		First(&tUsercompany).Error
	return tUsercompany, err
}

func (r repository) GetCompanyDefaultByUserId(conn *gorm.DB, userID string) (tUsercompany model.Usercompany, err error) {
	err = conn.Where("user_id = ? ", userID).
		Where("is_default_company = ? ", true).
		First(&tUsercompany).Error
	return tUsercompany, err
}

func (r repository) GetViewById(conn *gorm.DB, id string) (vUsercompany model.UsercompanyView, err error) {
	err = conn.Where("id = ? ", id).First(&vUsercompany).Error
	return vUsercompany, err
}

func (r repository) GetViewByUserIdAndCompanyId(conn *gorm.DB, userID, companyID string) (vUsercompany model.UsercompanyView, err error) {
	err = conn.Where("user_id = ? ", userID).
		Where("company_id = ? ", companyID).
		First(&vUsercompany).Error
	return vUsercompany, err
}

func (r repository) GetViewCreatorByCompanyId(conn *gorm.DB, companyID string) (vUsercompany model.UsercompanyView, err error) {
	err = conn.Where("company_id = ? ", companyID).
		Where("is_creator = ? ", true).
		First(&vUsercompany).Error
	return vUsercompany, err
}

func (r repository) GetViewCompanyDefaultByUserId(conn *gorm.DB, userID string) (vUsercompany model.UsercompanyView, err error) {
	err = conn.Where("user_id = ? ", userID).
		Where("is_default_company = ? ", true).
		First(&vUsercompany).Error
	return vUsercompany, err
}

func (r repository) Create(conn *gorm.DB, tUsercompany model.Usercompany) error {
	return conn.Create(&tUsercompany).Error
}

func (r repository) Update(conn *gorm.DB, tUsercompany model.Usercompany) error {
	return conn.Model(&tUsercompany).Updates(&tUsercompany).Error
}

func (r repository) Save(conn *gorm.DB, tUsercompany model.Usercompany) error {
	return conn.Save(&tUsercompany).Error
}

func (r repository) Delete(conn *gorm.DB, tUsercompany model.Usercompany) error {
	return conn.Delete(&tUsercompany).Error
}

func (r repository) Page(conn *gorm.DB, req request.PageUsercompany) (vUsercompanies []model.UsercompanyView, count int64, err error) {
	query := conn.Model(&vUsercompanies)
	if req.CompanyID != "" {
		query = query.Where("company_id = ?", req.CompanyID)
	}
	if req.UserID != "" {
		query = query.Where("user_id = ?", req.UserID)
	}

	err = query.Count(&count).Error
	if err != nil {
		return vUsercompanies, count, err
	}

	if req.SortField != "" {
		query = query.Order(fmt.Sprintf("%s %s", req.SortField, req.SortOrder))
	} else {
		query = query.Order(fmt.Sprintf("%s %s", "name", "asc"))
	}

	if req.Limit >= 0 {
		query = query.Offset((req.GetPage() - 1) * req.GetLimit()).Limit(req.GetLimit())
	}

	err = query.Find(&vUsercompanies).Error
	if err != nil {
		return vUsercompanies, count, err
	}

	return vUsercompanies, count, err
}

func NewRepository() Repository {
	return repository{}
}
