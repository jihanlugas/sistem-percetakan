package model

import (
	"github.com/jihanlugas/sistem-percetakan/utils"
	"gorm.io/gorm"
	"time"
)

var PreloadOrder = []string{
	"Company",
	"Customer",
	"Prints",
	"Prints.Paper",
	"Finishings",
	"Orderphases",
	"Transactions",
}

func (m *Order) BeforeCreate(tx *gorm.DB) (err error) {
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

func (m *Order) BeforeUpdate(tx *gorm.DB) (err error) {
	now := time.Now()
	m.UpdateDt = now
	return
}

//func (m *Order) BeforeDelete(tx *gorm.DB) (err error) {
//	return tx.Save(m).Error
//}
