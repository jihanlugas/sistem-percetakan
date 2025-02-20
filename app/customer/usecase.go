package customer

import (
	"errors"
	"fmt"
	"github.com/jihanlugas/sistem-percetakan/db"
	"github.com/jihanlugas/sistem-percetakan/jwt"
	"github.com/jihanlugas/sistem-percetakan/model"
	"github.com/jihanlugas/sistem-percetakan/request"
	"github.com/jihanlugas/sistem-percetakan/response"
	"github.com/jihanlugas/sistem-percetakan/utils"
)

type Usecase interface {
	Page(loginUser jwt.UserLogin, req request.PageCustomer) (vCustomers []model.CustomerView, count int64, err error)
	GetById(loginUser jwt.UserLogin, id string, preloads ...string) (vCustomer model.CustomerView, err error)
	Create(loginUser jwt.UserLogin, req request.CreateCustomer) error
	Update(loginUser jwt.UserLogin, id string, req request.UpdateCustomer) error
	Delete(loginUser jwt.UserLogin, id string) error
}

type usecase struct {
	repository Repository
}

func (u usecase) Page(loginUser jwt.UserLogin, req request.PageCustomer) (vCustomers []model.CustomerView, count int64, err error) {
	conn, closeConn := db.GetConnection()
	defer closeConn()

	vCustomers, count, err = u.repository.Page(conn, req)
	if err != nil {
		return vCustomers, count, err
	}

	return vCustomers, count, err
}

func (u usecase) GetById(loginUser jwt.UserLogin, id string, preloads ...string) (vCustomer model.CustomerView, err error) {
	conn, closeConn := db.GetConnection()
	defer closeConn()

	vCustomer, err = u.repository.GetViewById(conn, id, preloads...)
	if err != nil {
		return vCustomer, errors.New(fmt.Sprint("failed to get customer: ", err))
	}

	if jwt.IsSaveCompanyIDOR(loginUser, vCustomer.CompanyID) {
		return vCustomer, errors.New(response.ErrorHandlerIDOR)
	}

	return vCustomer, err
}

func (u usecase) Create(loginUser jwt.UserLogin, req request.CreateCustomer) error {
	var err error
	var tCustomer model.Customer

	conn, closeConn := db.GetConnection()
	defer closeConn()

	tx := conn.Begin()

	tCustomer = model.Customer{
		ID:          utils.GetUniqueID(),
		CompanyID:   req.CompanyID,
		Name:        req.Name,
		Description: req.Description,
		Email:       req.Email,
		PhoneNumber: utils.FormatPhoneTo62(req.PhoneNumber),
		Address:     req.Address,
		CreateBy:    loginUser.UserID,
		UpdateBy:    loginUser.UserID,
	}

	err = u.repository.Create(tx, tCustomer)
	if err != nil {
		return errors.New(fmt.Sprint("failed to create customer: ", err))
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return err
}

func (u usecase) Update(loginUser jwt.UserLogin, id string, req request.UpdateCustomer) error {
	var err error
	var tCustomer model.Customer

	conn, closeConn := db.GetConnection()
	defer closeConn()

	tCustomer, err = u.repository.GetTableById(conn, id)
	if err != nil {
		return errors.New(fmt.Sprint("failed to get customer: ", err))
	}

	if jwt.IsSaveCompanyIDOR(loginUser, tCustomer.CompanyID) {
		return errors.New(response.ErrorHandlerIDOR)
	}

	tx := conn.Begin()

	tCustomer.Name = req.Name
	tCustomer.Description = req.Description
	tCustomer.Address = req.Address
	tCustomer.Email = req.Email
	tCustomer.PhoneNumber = utils.FormatPhoneTo62(req.PhoneNumber)
	tCustomer.UpdateBy = loginUser.UserID
	err = u.repository.Save(tx, tCustomer)
	if err != nil {
		return errors.New(fmt.Sprint("failed to update customer: ", err))
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return err
}

func (u usecase) Delete(loginUser jwt.UserLogin, id string) error {
	var err error
	var tCustomer model.Customer

	conn, closeConn := db.GetConnection()
	defer closeConn()

	tCustomer, err = u.repository.GetTableById(conn, id)
	if err != nil {
		return errors.New(fmt.Sprint("failed to get customer: ", err))
	}

	if jwt.IsSaveCompanyIDOR(loginUser, tCustomer.CompanyID) {
		return errors.New(response.ErrorHandlerIDOR)
	}

	tx := conn.Begin()

	err = u.repository.Delete(tx, tCustomer)
	if err != nil {
		return errors.New(fmt.Sprint("failed to delete customer: ", err))
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return err
}

func NewUsecase(repository Repository) Usecase {
	return &usecase{
		repository: repository,
	}
}
