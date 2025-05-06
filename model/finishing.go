package model

import (
	"github.com/jihanlugas/sistem-percetakan/utils"
	"gorm.io/gorm"
	"time"
)

var PreloadFinishing = []string{
	"Company",
	"Order",
}

func (m *Finishing) BeforeCreate(tx *gorm.DB) (err error) {
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

func (m *Finishing) BeforeUpdate(tx *gorm.DB) (err error) {
	now := time.Now()
	m.UpdateDt = now
	return
}

//func (m *Finishing) BeforeDelete(tx *gorm.DB) (err error) {
//	return tx.Save(m).Error
//}
