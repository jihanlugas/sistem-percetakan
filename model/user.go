package model

import (
	"time"

	"github.com/jihanlugas/sistem-percetakan/utils"
	"gorm.io/gorm"
)

func (m *User) BeforeCreate(tx *gorm.DB) (err error) {
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

func (m *User) BeforeUpdate(tx *gorm.DB) (err error) {
	now := time.Now()
	m.UpdateDt = now
	return
}

func (m *UserView) AfterFind(tx *gorm.DB) (err error) {
	if m.PhotoID != "" {
		// m.PhotoUrl = fmt.Sprintf("%s/%s", config.FileBaseUrl, m.PhotoUrl)
		m.PhotoUrl = utils.GetPhotoUrlById(m.PhotoID)
	}
	return
}
