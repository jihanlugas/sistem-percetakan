package model

import (
	"github.com/jihanlugas/sistem-percetakan/utils"
	"gorm.io/gorm"
	"time"
)

func (m *Paper) BeforeCreate(tx *gorm.DB) (err error) {
	now := time.Now()

	if m.ID == "" {
		m.ID = utils.GetUniqueID()
	}

	m.CreateDt = now
	m.UpdateDt = now
	return
}

func (m *Paper) BeforeUpdate(tx *gorm.DB) (err error) {
	now := time.Now()
	m.UpdateDt = now
	return
}

func (m *Paper) BeforeDelete(tx *gorm.DB) (err error) {
	return tx.Save(m).Error
}
