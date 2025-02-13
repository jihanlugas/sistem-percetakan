package payment

import (
	"errors"
	"fmt"
	"github.com/jihanlugas/sistem-percetakan/app/auth"
	"github.com/jihanlugas/sistem-percetakan/db"
	"github.com/jihanlugas/sistem-percetakan/jwt"
	"github.com/jihanlugas/sistem-percetakan/model"
	"github.com/jihanlugas/sistem-percetakan/request"
	"github.com/jihanlugas/sistem-percetakan/response"
	"github.com/jihanlugas/sistem-percetakan/utils"
)

type Usecase interface {
	Page(loginUser jwt.UserLogin, req request.PagePayment) (vPayments []model.PaymentView, count int64, err error)
	GetById(loginUser jwt.UserLogin, id string, preloads ...string) (vPayment model.PaymentView, err error)
	Create(loginUser jwt.UserLogin, req request.CreatePayment) error
	Update(loginUser jwt.UserLogin, id string, req request.UpdatePayment) error
	Delete(loginUser jwt.UserLogin, id string) error
}

type usecase struct {
	repository Repository
}

func (u usecase) Page(loginUser jwt.UserLogin, req request.PagePayment) (vPayments []model.PaymentView, count int64, err error) {
	conn, closeConn := db.GetConnection()
	defer closeConn()

	vPayments, count, err = u.repository.Page(conn, req)
	if err != nil {
		return vPayments, count, err
	}

	return vPayments, count, err
}

func (u usecase) GetById(loginUser jwt.UserLogin, id string, preloads ...string) (vPayment model.PaymentView, err error) {
	conn, closeConn := db.GetConnection()
	defer closeConn()

	vPayment, err = u.repository.GetViewById(conn, id, preloads...)
	if err != nil {
		return vPayment, errors.New(fmt.Sprint("failed to get order: ", err))
	}

	if auth.IsSaveIDOR(loginUser, vPayment.CompanyID) {
		return vPayment, errors.New(response.ErrorHandlerIDOR)
	}

	return vPayment, err
}

func (u usecase) Create(loginUser jwt.UserLogin, req request.CreatePayment) error {
	var err error
	var tPayment model.Payment

	conn, closeConn := db.GetConnection()
	defer closeConn()

	if auth.IsSaveIDOR(loginUser, req.CompanyID) {
		return errors.New(response.ErrorHandlerIDOR)
	}

	tx := conn.Begin()

	tPayment = model.Payment{
		ID:          utils.GetUniqueID(),
		CompanyID:   req.CompanyID,
		OrderID:     req.OrderID,
		Name:        req.Name,
		Description: req.Description,
		Amount:      req.Amount,
		CreateBy:    loginUser.UserID,
		UpdateBy:    loginUser.UserID,
	}

	err = u.repository.Create(tx, tPayment)
	if err != nil {
		return errors.New(fmt.Sprint("failed to create payment: ", err))
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return err
}

func (u usecase) Update(loginUser jwt.UserLogin, id string, req request.UpdatePayment) error {
	var err error
	var tPayment model.Payment

	conn, closeConn := db.GetConnection()
	defer closeConn()

	tPayment, err = u.repository.GetTableById(conn, id)
	if err != nil {
		return errors.New(fmt.Sprint("failed to get payment: ", err))
	}

	if auth.IsSaveIDOR(loginUser, tPayment.CompanyID) {
		return errors.New(response.ErrorHandlerIDOR)
	}

	tx := conn.Begin()

	tPayment.Name = req.Name
	tPayment.Description = req.Description
	tPayment.Amount = req.Amount
	tPayment.UpdateBy = loginUser.UserID
	err = u.repository.Save(tx, tPayment)
	if err != nil {
		return errors.New(fmt.Sprint("failed to update payment: ", err))
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return err
}

func (u usecase) Delete(loginUser jwt.UserLogin, id string) error {
	var err error
	var tPayment model.Payment

	conn, closeConn := db.GetConnection()
	defer closeConn()

	tPayment, err = u.repository.GetTableById(conn, id)
	if err != nil {
		return errors.New(fmt.Sprint("failed to get payment: ", err))
	}

	if auth.IsSaveIDOR(loginUser, tPayment.CompanyID) {
		return errors.New(response.ErrorHandlerIDOR)
	}

	tx := conn.Begin()

	err = u.repository.Delete(tx, tPayment)
	if err != nil {
		return errors.New(fmt.Sprint("failed to delete payment: ", err))
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
