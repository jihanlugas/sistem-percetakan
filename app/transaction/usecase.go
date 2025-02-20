package transaction

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
	Page(loginUser jwt.UserLogin, req request.PageTransaction) (vTransactions []model.TransactionView, count int64, err error)
	GetById(loginUser jwt.UserLogin, id string, preloads ...string) (vTransaction model.TransactionView, err error)
	Create(loginUser jwt.UserLogin, req request.CreateTransaction) error
	Update(loginUser jwt.UserLogin, id string, req request.UpdateTransaction) error
	Delete(loginUser jwt.UserLogin, id string) error
}

type usecase struct {
	repository Repository
}

func (u usecase) Page(loginUser jwt.UserLogin, req request.PageTransaction) (vTransactions []model.TransactionView, count int64, err error) {
	conn, closeConn := db.GetConnection()
	defer closeConn()

	vTransactions, count, err = u.repository.Page(conn, req)
	if err != nil {
		return vTransactions, count, err
	}

	return vTransactions, count, err
}

func (u usecase) GetById(loginUser jwt.UserLogin, id string, preloads ...string) (vTransaction model.TransactionView, err error) {
	conn, closeConn := db.GetConnection()
	defer closeConn()

	vTransaction, err = u.repository.GetViewById(conn, id, preloads...)
	if err != nil {
		return vTransaction, errors.New(fmt.Sprint("failed to get order: ", err))
	}

	if jwt.IsSaveCompanyIDOR(loginUser, vTransaction.CompanyID) {
		return vTransaction, errors.New(response.ErrorHandlerIDOR)
	}

	return vTransaction, err
}

func (u usecase) Create(loginUser jwt.UserLogin, req request.CreateTransaction) error {
	var err error
	var tTransaction model.Transaction

	conn, closeConn := db.GetConnection()
	defer closeConn()

	if jwt.IsSaveCompanyIDOR(loginUser, req.CompanyID) {
		return errors.New(response.ErrorHandlerIDOR)
	}

	tx := conn.Begin()

	tTransaction = model.Transaction{
		ID:          utils.GetUniqueID(),
		CompanyID:   req.CompanyID,
		OrderID:     req.OrderID,
		Name:        req.Name,
		Description: req.Description,
		Amount:      req.Amount,
		Type:        req.Type,
		CreateBy:    loginUser.UserID,
		UpdateBy:    loginUser.UserID,
	}

	err = u.repository.Create(tx, tTransaction)
	if err != nil {
		return errors.New(fmt.Sprint("failed to create transaction: ", err))
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return err
}

func (u usecase) Update(loginUser jwt.UserLogin, id string, req request.UpdateTransaction) error {
	var err error
	var tTransaction model.Transaction

	conn, closeConn := db.GetConnection()
	defer closeConn()

	tTransaction, err = u.repository.GetTableById(conn, id)
	if err != nil {
		return errors.New(fmt.Sprint("failed to get transaction: ", err))
	}

	if jwt.IsSaveCompanyIDOR(loginUser, tTransaction.CompanyID) {
		return errors.New(response.ErrorHandlerIDOR)
	}

	tx := conn.Begin()

	tTransaction.Name = req.Name
	tTransaction.Description = req.Description
	tTransaction.Amount = req.Amount
	tTransaction.Type = req.Type
	tTransaction.UpdateBy = loginUser.UserID
	err = u.repository.Save(tx, tTransaction)
	if err != nil {
		return errors.New(fmt.Sprint("failed to update transaction: ", err))
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return err
}

func (u usecase) Delete(loginUser jwt.UserLogin, id string) error {
	var err error
	var tTransaction model.Transaction

	conn, closeConn := db.GetConnection()
	defer closeConn()

	tTransaction, err = u.repository.GetTableById(conn, id)
	if err != nil {
		return errors.New(fmt.Sprint("failed to get transaction: ", err))
	}

	if jwt.IsSaveCompanyIDOR(loginUser, tTransaction.CompanyID) {
		return errors.New(response.ErrorHandlerIDOR)
	}

	tx := conn.Begin()

	err = u.repository.Delete(tx, tTransaction)
	if err != nil {
		return errors.New(fmt.Sprint("failed to delete transaction: ", err))
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
