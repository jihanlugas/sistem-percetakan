package model

import (
	"fmt"
	"github.com/jihanlugas/sistem-percetakan/config"
	"github.com/jihanlugas/sistem-percetakan/utils"
	"gorm.io/gorm"
	"time"
)

func (m *User) BeforeCreate(tx *gorm.DB) (err error) {
	now := time.Now()

	if m.ID == "" {
		m.ID = utils.GetUniqueID()
	}

	m.CreateDt = now
	m.UpdateDt = now
	return
}

func (m *User) BeforeUpdate(tx *gorm.DB) (err error) {
	now := time.Now()
	m.UpdateDt = now
	return
}

func (m *User) BeforeDelete(tx *gorm.DB) (err error) {
	return tx.Save(m).Error
}

func (m *UserView) AfterFind(tx *gorm.DB) (err error) {
	if m.PhotoUrl != "" {
		m.PhotoUrl = fmt.Sprintf("%s/%s", config.FileBaseUrl, m.PhotoUrl)
	}
	return
}
