package orderphase

import (
	"fmt"
	"github.com/jihanlugas/sistem-percetakan/model"
	"github.com/jihanlugas/sistem-percetakan/request"
	"gorm.io/gorm"
)

type Repository interface {
	GetTableById(conn *gorm.DB, id string, preloads ...string) (tOrderphase model.Orderphase, err error)
	GetViewById(conn *gorm.DB, id string, preloads ...string) (vOrderphase model.OrderphaseView, err error)
	Create(conn *gorm.DB, tOrderphase model.Orderphase) error
	Creates(conn *gorm.DB, tOrderphases []model.Orderphase) error
	Update(conn *gorm.DB, tOrderphase model.Orderphase) error
	Save(conn *gorm.DB, tOrderphase model.Orderphase) error
	Delete(conn *gorm.DB, tOrderphase model.Orderphase) error
	Page(conn *gorm.DB, req request.PageOrderphase) (vOrderphases []model.OrderphaseView, count int64, err error)
}

type repository struct {
}

func (r repository) GetTableById(conn *gorm.DB, id string, preloads ...string) (tOrderphase model.Orderphase, err error) {
	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}
	err = conn.Where("id = ? ", id).First(&tOrderphase).Error
	return tOrderphase, err
}

func (r repository) GetViewById(conn *gorm.DB, id string, preloads ...string) (vOrderphase model.OrderphaseView, err error) {
	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}
	err = conn.Where("id = ? ", id).First(&vOrderphase).Error
	return vOrderphase, err
}

func (r repository) Create(conn *gorm.DB, tOrderphase model.Orderphase) error {
	return conn.Create(&tOrderphase).Error
}

func (r repository) Creates(conn *gorm.DB, tOrderphases []model.Orderphase) error {
	return conn.Create(&tOrderphases).Error
}

func (r repository) Update(conn *gorm.DB, tOrderphase model.Orderphase) error {
	return conn.Model(&tOrderphase).Updates(&tOrderphase).Error
}

func (r repository) Save(conn *gorm.DB, tOrderphase model.Orderphase) error {
	return conn.Save(&tOrderphase).Error
}

func (r repository) Delete(conn *gorm.DB, tOrderphase model.Orderphase) error {
	return conn.Delete(&tOrderphase).Error
}

func (r repository) Page(conn *gorm.DB, req request.PageOrderphase) (vOrderphases []model.OrderphaseView, count int64, err error) {
	query := conn.Model(&vOrderphases)

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
		return vOrderphases, count, err
	}

	if req.SortField != "" {
		query = query.Order(fmt.Sprintf("%s %s", req.SortField, req.SortOrder))
	} else {
		query = query.Order(fmt.Sprintf("%s %s", "name", "asc"))
	}
	err = query.Offset((req.GetPage() - 1) * req.GetLimit()).
		Limit(req.GetLimit()).
		Find(&vOrderphases).Error
	if err != nil {
		return vOrderphases, count, err
	}

	return vOrderphases, count, err
}

func NewRepository() Repository {
	return &repository{}
}
