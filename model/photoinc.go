package model

import (
	"github.com/jihanlugas/sistem-percetakan/utils"
	"gorm.io/gorm"
)

func (m *Photoinc) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = utils.GetUniqueID()
	}

	return
}
