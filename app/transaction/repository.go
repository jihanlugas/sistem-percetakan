package transaction

import (
	"fmt"
	"github.com/jihanlugas/sistem-percetakan/constant"
	"github.com/jihanlugas/sistem-percetakan/model"
	"github.com/jihanlugas/sistem-percetakan/request"
	"github.com/jihanlugas/sistem-percetakan/utils"
	"gorm.io/gorm"
	"strings"
	"time"
)

type Repository interface {
	GetTableById(conn *gorm.DB, id string, preloads ...string) (tTransaction model.Transaction, err error)
	GetViewById(conn *gorm.DB, id string, preloads ...string) (vTransaction model.TransactionView, err error)
	Create(conn *gorm.DB, tTransaction model.Transaction) error
	Creates(conn *gorm.DB, tTransactions []model.Transaction) error
	Update(conn *gorm.DB, tTransaction model.Transaction) error
	Save(conn *gorm.DB, tTransaction model.Transaction) error
	Delete(conn *gorm.DB, tTransaction model.Transaction) error
	DeleteByOrderId(conn *gorm.DB, id string) error
	Page(conn *gorm.DB, req request.PageTransaction) (vTransactions []model.TransactionView, count int64, err error)
	GetDailyAmountPeriod(conn *gorm.DB, companyID string, transactionType constant.TransactionType, startDt, endDt time.Time) (data []model.TransactionPeriod, err error)
	GetTotalAmountPeriod(conn *gorm.DB, companyID string, transactionType constant.TransactionType, startDt, endDt time.Time) (total_amount int64, err error)
}

type repository struct {
}

func (r repository) GetTableById(conn *gorm.DB, id string, preloads ...string) (tTransaction model.Transaction, err error) {
	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}
	err = conn.Where("id = ? ", id).First(&tTransaction).Error
	return tTransaction, err
}

func (r repository) GetViewById(conn *gorm.DB, id string, preloads ...string) (vTransaction model.TransactionView, err error) {
	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}
	err = conn.Where("id = ? ", id).First(&vTransaction).Error
	return vTransaction, err
}

func (r repository) Create(conn *gorm.DB, tTransaction model.Transaction) error {
	return conn.Create(&tTransaction).Error
}

func (r repository) Creates(conn *gorm.DB, tTransactions []model.Transaction) error {
	return conn.Create(&tTransactions).Error
}

func (r repository) Update(conn *gorm.DB, tTransaction model.Transaction) error {
	return conn.Model(&tTransaction).Updates(&tTransaction).Error
}

func (r repository) Save(conn *gorm.DB, tTransaction model.Transaction) error {
	return conn.Save(&tTransaction).Error
}

func (r repository) Delete(conn *gorm.DB, tTransaction model.Transaction) error {
	return conn.Delete(&tTransaction).Error
}

func (r repository) DeleteByOrderId(conn *gorm.DB, id string) error {
	return conn.Where("order_id = ? ", id).Delete(&model.Transaction{}).Error
}

func (r repository) Page(conn *gorm.DB, req request.PageTransaction) (vTransactions []model.TransactionView, count int64, err error) {
	query := conn.Model(&vTransactions)

	// preloads
	preloads := strings.Split(req.Preloads, ",")
	for _, preload := range preloads {
		if utils.IsAvailablePreload(preload, model.PreloadTransaction) {
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
	if req.Type != 0 {
		query = query.Where("type = ?", req.Type)
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
	if req.StartAmount != nil {
		query = query.Where("amount >= ?", req.StartAmount)
	}
	if req.EndAmount != nil {
		query = query.Where("amount <= ?", req.EndAmount)
	}

	err = query.Count(&count).Error
	if err != nil {
		return vTransactions, count, err
	}

	if req.SortField != "" {
		query = query.Order(fmt.Sprintf("%s %s", req.SortField, req.SortOrder))
	} else {
		query = query.Order(fmt.Sprintf("%s %s", "name", "asc"))
	}

	if req.Limit >= 0 {
		query = query.Offset((req.GetPage() - 1) * req.GetLimit()).Limit(req.GetLimit())
	}

	err = query.Find(&vTransactions).Error
	if err != nil {
		return vTransactions, count, err
	}

	return vTransactions, count, err
}

func (r repository) GetDailyAmountPeriod(conn *gorm.DB, companyID string, transactionType constant.TransactionType, startDt, endDt time.Time) (data []model.TransactionPeriod, err error) {
	err = conn.Raw(`WITH date_series AS (
			SELECT generate_series(
				?::DATE, 
				?::DATE, 
				'1 day'
			) AS date
		)
		SELECT 
			ds.date, 
			COALESCE(SUM(t.amount), 0) AS amount
		FROM date_series ds
		LEFT JOIN transactions t 
			ON DATE(t.create_dt) = ds.date
			and t.company_id = ?
			and t."type" = ?
		GROUP BY ds.date
		ORDER BY ds.date ASC`, startDt, endDt, companyID, transactionType).Scan(&data).Error
	fmt.Print(data)
	return data, err
}

func (r repository) GetTotalAmountPeriod(conn *gorm.DB, companyID string, transactionType constant.TransactionType, startDt, endDt time.Time) (total_amount int64, err error) {
	err = conn.Raw(`
		SELECT COALESCE(SUM(amount), 0) AS total_amount
		FROM transactions
		WHERE company_id = ?
		AND type = ?
		AND create_dt BETWEEN ? AND ?
		AND delete_dt is null
	`, companyID, transactionType, startDt, endDt).Scan(&total_amount).Error

	return total_amount, err
}

func NewRepository() Repository {
	return &repository{}
}
