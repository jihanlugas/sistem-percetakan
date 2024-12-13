package user

import (
	"fmt"
	"github.com/jihanlugas/sistem-percetakan/model"
	"github.com/jihanlugas/sistem-percetakan/request"
	"github.com/jihanlugas/sistem-percetakan/utils"
	"gorm.io/gorm"
)

type Repository interface {
	GetById(conn *gorm.DB, id string) (tUser model.User, err error)
	GetByUsername(conn *gorm.DB, username string) (tUser model.User, err error)
	GetByEmail(conn *gorm.DB, email string) (tUser model.User, err error)
	GetByNoHp(conn *gorm.DB, noHp string) (tUser model.User, err error)
	GetViewById(conn *gorm.DB, id string) (vUser model.UserView, err error)
	GetViewByUsername(conn *gorm.DB, username string) (vUser model.UserView, err error)
	GetViewByEmail(conn *gorm.DB, email string) (vUser model.UserView, err error)
	GetViewByNoHp(conn *gorm.DB, noHp string) (vUser model.UserView, err error)
	Create(conn *gorm.DB, tUser model.User) error
	Update(conn *gorm.DB, tUser model.User) error
	Save(conn *gorm.DB, tUser model.User) error
	Delete(conn *gorm.DB, tUser model.User) error
	Page(conn *gorm.DB, req request.PageUser) (vUsers []model.UserView, count int64, err error)
}

type repository struct {
}

func (r repository) GetById(conn *gorm.DB, id string) (tUser model.User, err error) {
	err = conn.Where("id = ? ", id).First(&tUser).Error
	return tUser, err
}

func (r repository) GetByUsername(conn *gorm.DB, username string) (tUser model.User, err error) {
	err = conn.Where("username = ? ", username).First(&tUser).Error
	return tUser, err
}

func (r repository) GetByEmail(conn *gorm.DB, email string) (tUser model.User, err error) {
	err = conn.Where("email = ? ", email).First(&tUser).Error
	return tUser, err
}

func (r repository) GetByNoHp(conn *gorm.DB, noHp string) (tUser model.User, err error) {
	err = conn.Where("no_hp = ? ", utils.FormatPhoneTo62(noHp)).First(&tUser).Error
	return tUser, err
}

func (r repository) GetViewById(conn *gorm.DB, id string) (vUser model.UserView, err error) {
	err = conn.Where("id = ? ", id).First(&vUser).Error
	return vUser, err
}

func (r repository) GetViewByUsername(conn *gorm.DB, username string) (vUser model.UserView, err error) {
	err = conn.Where("username = ? ", username).First(&vUser).Error
	return vUser, err
}

func (r repository) GetViewByEmail(conn *gorm.DB, email string) (vUser model.UserView, err error) {
	err = conn.Where("email = ? ", email).First(&vUser).Error
	return vUser, err
}

func (r repository) GetViewByNoHp(conn *gorm.DB, noHp string) (vUser model.UserView, err error) {
	err = conn.Where("no_hp = ? ", noHp).First(&vUser).Error
	return vUser, err
}

func (r repository) Create(conn *gorm.DB, tUser model.User) error {
	return conn.Create(&tUser).Error
}

func (r repository) Update(conn *gorm.DB, tUser model.User) error {
	return conn.Model(&tUser).Updates(&tUser).Error
}

func (r repository) Save(conn *gorm.DB, tUser model.User) error {
	return conn.Save(&tUser).Error
}

func (r repository) Delete(conn *gorm.DB, tUser model.User) error {
	return conn.Delete(&tUser).Error
}

func (r repository) Page(conn *gorm.DB, req request.PageUser) (vUsers []model.UserView, count int64, err error) {
	query := conn.Model(&vUsers)
	if req.Email != "" {
		query = query.Where("email ILIKE ?", "%"+req.Email+"%")
	}
	if req.Username != "" {
		query = query.Where("username ILIKE ?", "%"+req.Username+"%")
	}
	if req.NoHp != "" {
		query = query.Where("no_hp ILIKE ?", "%"+utils.FormatPhoneTo62(req.NoHp)+"%")
	}
	if req.Fullname != "" {
		query = query.Where("fullname ILIKE ?", "%"+req.Fullname+"%")
	}

	err = query.Count(&count).Error
	if err != nil {
		return vUsers, count, err
	}

	if req.SortField != "" {
		query = query.Order(fmt.Sprintf("%s %s", req.SortField, req.SortOrder))
	} else {
		query = query.Order(fmt.Sprintf("%s %s", "fullname", "asc"))
	}
	err = query.Offset((req.GetPage() - 1) * req.GetLimit()).
		Limit(req.GetLimit()).
		Find(&vUsers).Error
	if err != nil {
		return vUsers, count, err
	}

	return vUsers, count, err
}

func NewRepository() Repository {
	return repository{}
}
