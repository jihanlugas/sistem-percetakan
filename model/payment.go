package model

import (
	"github.com/jihanlugas/sistem-percetakan/utils"
	"gorm.io/gorm"
	"time"
)

func (m *Payment) BeforeCreate(tx *gorm.DB) (err error) {
	now := time.Now()

	if m.ID == "" {
		m.ID = utils.GetUniqueID()
	}

	m.CreateDt = now
	m.UpdateDt = now
	return
}

func (m *Payment) BeforeUpdate(tx *gorm.DB) (err error) {
	now := time.Now()
	m.UpdateDt = now
	return
}

func (m *Payment) BeforeDelete(tx *gorm.DB) (err error) {
	return tx.Save(m).Error
}
