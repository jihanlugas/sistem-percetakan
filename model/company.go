package model

import (
	"github.com/jihanlugas/sistem-percetakan/utils"
	"gorm.io/gorm"
	"time"
)

func (m *Company) BeforeCreate(tx *gorm.DB) (err error) {
	now := time.Now()

	if m.ID == "" {
		m.ID = utils.GetUniqueID()
	}

	m.CreateDt = now
	m.UpdateDt = now
	return
}

func (m *Company) BeforeUpdate(tx *gorm.DB) (err error) {
	now := time.Now()
	m.UpdateDt = now
	return
}

func (m *Company) BeforeDelete(tx *gorm.DB) (err error) {
	return tx.Save(m).Error
}
