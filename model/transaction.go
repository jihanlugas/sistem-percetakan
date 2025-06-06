package model

import (
	"github.com/jihanlugas/sistem-percetakan/utils"
	"gorm.io/gorm"
	"time"
)

var PreloadTransaction = []string{
	"Company",
	"Order",
}

type TransactionPeriod struct {
	Date   time.Time `json:"date"`
	Amount int64     `json:"amount"`
}

func (m *Transaction) BeforeCreate(tx *gorm.DB) (err error) {
	now := time.Now()

	if m.ID == "" {
		m.ID = utils.GetUniqueID()
	}

	if m.CreateDt.IsZero() {
		m.CreateDt = now
	}
	if m.UpdateDt.IsZero() {
		m.UpdateDt = now
	}
	return
}

func (m *Transaction) BeforeUpdate(tx *gorm.DB) (err error) {
	now := time.Now()
	m.UpdateDt = now
	return
}

//func (m *Transaction) BeforeDelete(tx *gorm.DB) (err error) {
//	return tx.Save(m).Error
//}
