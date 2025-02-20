package user

import (
	"errors"
	"fmt"
	"github.com/jihanlugas/sistem-percetakan/app/usercompany"
	"github.com/jihanlugas/sistem-percetakan/constant"
	"github.com/jihanlugas/sistem-percetakan/cryption"
	"github.com/jihanlugas/sistem-percetakan/db"
	"github.com/jihanlugas/sistem-percetakan/jwt"
	"github.com/jihanlugas/sistem-percetakan/model"
	"github.com/jihanlugas/sistem-percetakan/request"
	"github.com/jihanlugas/sistem-percetakan/response"
	"github.com/jihanlugas/sistem-percetakan/utils"
	"time"
)

type Usecase interface {
	Page(loginUser jwt.UserLogin, req request.PageUser) (vUsers []model.UserView, count int64, err error)
	GetById(loginUser jwt.UserLogin, id string, preloads ...string) (vUser model.UserView, err error)
	Create(loginUser jwt.UserLogin, req request.CreateUser) error
	Update(loginUser jwt.UserLogin, id string, req request.UpdateUser) error
	ChangePassword(loginUser jwt.UserLogin, req request.ChangePassword) error
	Delete(loginUser jwt.UserLogin, id string) error
}

type usecase struct {
	userRepository        Repository
	usercompanyRepository usercompany.Repository
}

func (u usecase) Page(loginUser jwt.UserLogin, req request.PageUser) (vUsers []model.UserView, count int64, err error) {
	conn, closeConn := db.GetConnection()
	defer closeConn()

	vUsers, count, err = u.userRepository.Page(conn, req)
	if err != nil {
		return vUsers, count, err
	}

	return vUsers, count, err
}

func (u usecase) GetById(loginUser jwt.UserLogin, id string, preloads ...string) (vUser model.UserView, err error) {
	conn, closeConn := db.GetConnection()
	defer closeConn()

	vUser, err = u.userRepository.GetViewById(conn, id, preloads...)
	if err != nil {
		return vUser, errors.New(fmt.Sprint("failed to get user: ", err))
	}

	vUsercompany, err := u.usercompanyRepository.GetViewByUserIdAndCompanyId(conn, vUser.ID, loginUser.CompanyID)
	if err != nil {
		return vUser, errors.New(fmt.Sprint("failed to get user company: ", err))
	}

	if jwt.IsSaveCompanyIDOR(loginUser, vUsercompany.CompanyID) {
		return vUser, errors.New(response.ErrorHandlerIDOR)
	}

	return vUser, err
}

func (u usecase) Create(loginUser jwt.UserLogin, req request.CreateUser) error {
	var err error
	var tUser model.User

	conn, closeConn := db.GetConnection()
	defer closeConn()

	tx := conn.Begin()

	now := time.Now()

	encodePasswd, err := cryption.EncryptAES64(req.Passwd)
	if err != nil {
		return errors.New(fmt.Sprint("failed to encode password: ", err))
	}

	tUser = model.User{
		ID:                utils.GetUniqueID(),
		Role:              constant.RoleUser,
		Email:             req.Email,
		Username:          req.Username,
		PhoneNumber:       utils.FormatPhoneTo62(req.PhoneNumber),
		Address:           req.Address,
		Fullname:          req.Fullname,
		Passwd:            encodePasswd,
		PassVersion:       1,
		IsActive:          true,
		PhotoID:           "",
		LastLoginDt:       nil,
		BirthDt:           req.BirthDt,
		BirthPlace:        req.BirthPlace,
		AccountVerifiedDt: &now,
		CreateBy:          loginUser.UserID,
		UpdateBy:          loginUser.UserID,
	}

	err = u.userRepository.Create(tx, tUser)
	if err != nil {
		return errors.New(fmt.Sprint("failed to create user: ", err))
	}

	tUsercompany := model.Usercompany{
		ID:               utils.GetUniqueID(),
		UserID:           tUser.ID,
		CompanyID:        loginUser.CompanyID,
		IsDefaultCompany: true,
		IsCreator:        false,
		CreateBy:         loginUser.UserID,
		UpdateBy:         loginUser.UserID,
	}
	err = u.usercompanyRepository.Create(tx, tUsercompany)
	if err != nil {
		return errors.New(fmt.Sprint("failed to create user company: ", err))
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return err
}

func (u usecase) Update(loginUser jwt.UserLogin, id string, req request.UpdateUser) error {
	var err error
	var tUser model.User

	conn, closeConn := db.GetConnection()
	defer closeConn()

	tUser, err = u.userRepository.GetTableById(conn, id)
	if err != nil {
		return errors.New(fmt.Sprint("failed to get user: ", err))
	}

	vUsercompany, err := u.usercompanyRepository.GetViewByUserIdAndCompanyId(conn, tUser.ID, loginUser.CompanyID)
	if err != nil {
		return errors.New(fmt.Sprint("failed to get user company: ", err))
	}

	if jwt.IsSaveCompanyIDOR(loginUser, vUsercompany.CompanyID) {
		return errors.New(response.ErrorHandlerIDOR)
	}

	tx := conn.Begin()

	tUser.Fullname = req.Fullname
	tUser.Email = req.Email
	tUser.PhoneNumber = utils.FormatPhoneTo62(req.PhoneNumber)
	tUser.Username = req.Username
	tUser.Address = req.Address
	tUser.BirthDt = req.BirthDt
	tUser.BirthPlace = req.BirthPlace
	tUser.UpdateBy = loginUser.UserID
	err = u.userRepository.Save(tx, tUser)
	if err != nil {
		return errors.New(fmt.Sprint("failed to update user: ", err))
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return err
}

func (u usecase) ChangePassword(loginUser jwt.UserLogin, req request.ChangePassword) error {
	var err error
	var tUser model.User

	conn, closeConn := db.GetConnection()
	defer closeConn()

	tUser, err = u.userRepository.GetTableById(conn, loginUser.UserID)
	if err != nil {
		return errors.New(fmt.Sprint("failed to get user: ", err))
	}

	tx := conn.Begin()

	err = cryption.CheckAES64(req.CurrentPasswd, tUser.Passwd)
	if err != nil {
		return errors.New(fmt.Sprint("invalid current password"))
	}

	encodePasswd, err := cryption.EncryptAES64(req.Passwd)
	if err != nil {
		return errors.New(fmt.Sprint("failed to encode password: ", err))
	}

	tUser.Passwd = encodePasswd
	tUser.PassVersion += 1
	tUser.UpdateBy = loginUser.UserID
	err = u.userRepository.Save(tx, tUser)
	if err != nil {
		return errors.New(fmt.Sprint("failed to update password: ", err))
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return err
}

func (u usecase) Delete(loginUser jwt.UserLogin, id string) error {
	var err error
	var tUser model.User

	conn, closeConn := db.GetConnection()
	defer closeConn()

	tUser, err = u.userRepository.GetTableById(conn, id)
	if err != nil {
		return errors.New(fmt.Sprint("failed to get user: ", err))
	}

	vUsercompany, err := u.usercompanyRepository.GetViewByUserIdAndCompanyId(conn, tUser.ID, loginUser.CompanyID)
	if err != nil {
		return errors.New(fmt.Sprint("failed to get user company: ", err))
	}

	if jwt.IsSaveCompanyIDOR(loginUser, vUsercompany.CompanyID) {
		return errors.New(response.ErrorHandlerIDOR)
	}

	tx := conn.Begin()

	err = u.userRepository.Delete(tx, tUser)
	if err != nil {
		return errors.New(fmt.Sprint("failed to delete user: ", err))
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return err
}

func NewUsecase(userRepository Repository, usercompanyRepository usercompany.Repository) Usecase {
	return &usecase{
		userRepository:        userRepository,
		usercompanyRepository: usercompanyRepository,
	}
}
