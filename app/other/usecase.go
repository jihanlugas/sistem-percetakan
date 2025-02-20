package other

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
	Page(loginUser jwt.UserLogin, req request.PageOther) (vOthers []model.OtherView, count int64, err error)
	GetById(loginUser jwt.UserLogin, id string, preloads ...string) (vOther model.OtherView, err error)
	Create(loginUser jwt.UserLogin, req request.CreateOther) error
	Update(loginUser jwt.UserLogin, id string, req request.UpdateOther) error
	Delete(loginUser jwt.UserLogin, id string) error
}

type usecase struct {
	repository Repository
}

func (u usecase) Page(loginUser jwt.UserLogin, req request.PageOther) (vOthers []model.OtherView, count int64, err error) {
	conn, closeConn := db.GetConnection()
	defer closeConn()

	vOthers, count, err = u.repository.Page(conn, req)
	if err != nil {
		return vOthers, count, err
	}

	return vOthers, count, err
}

func (u usecase) GetById(loginUser jwt.UserLogin, id string, preloads ...string) (vOther model.OtherView, err error) {
	conn, closeConn := db.GetConnection()
	defer closeConn()

	vOther, err = u.repository.GetViewById(conn, id, preloads...)
	if err != nil {
		return vOther, errors.New(fmt.Sprint("failed to get order: ", err))
	}

	if jwt.IsSaveCompanyIDOR(loginUser, vOther.CompanyID) {
		return vOther, errors.New(response.ErrorHandlerIDOR)
	}

	return vOther, err
}

func (u usecase) Create(loginUser jwt.UserLogin, req request.CreateOther) error {
	var err error
	var tOther model.Other

	conn, closeConn := db.GetConnection()
	defer closeConn()

	if jwt.IsSaveCompanyIDOR(loginUser, req.CompanyID) {
		return errors.New(response.ErrorHandlerIDOR)
	}

	tx := conn.Begin()

	tOther = model.Other{
		ID:          utils.GetUniqueID(),
		CompanyID:   req.CompanyID,
		OrderID:     req.OrderID,
		Name:        req.Name,
		Description: req.Description,
		Qty:         req.Qty,
		Price:       req.Price,
		Total:       req.Total,
		CreateBy:    loginUser.UserID,
		UpdateBy:    loginUser.UserID,
	}

	err = u.repository.Create(tx, tOther)
	if err != nil {
		return errors.New(fmt.Sprint("failed to create other: ", err))
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return err
}

func (u usecase) Update(loginUser jwt.UserLogin, id string, req request.UpdateOther) error {
	var err error
	var tOther model.Other

	conn, closeConn := db.GetConnection()
	defer closeConn()

	tOther, err = u.repository.GetTableById(conn, id)
	if err != nil {
		return errors.New(fmt.Sprint("failed to get other: ", err))
	}

	if jwt.IsSaveCompanyIDOR(loginUser, tOther.CompanyID) {
		return errors.New(response.ErrorHandlerIDOR)
	}

	tx := conn.Begin()

	tOther.Name = req.Name
	tOther.Description = req.Description
	tOther.Qty = req.Qty
	tOther.Price = req.Price
	tOther.Total = req.Total
	tOther.UpdateBy = loginUser.UserID
	err = u.repository.Save(tx, tOther)
	if err != nil {
		return errors.New(fmt.Sprint("failed to update other: ", err))
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return err
}

func (u usecase) Delete(loginUser jwt.UserLogin, id string) error {
	var err error
	var tOther model.Other

	conn, closeConn := db.GetConnection()
	defer closeConn()

	tOther, err = u.repository.GetTableById(conn, id)
	if err != nil {
		return errors.New(fmt.Sprint("failed to get other: ", err))
	}

	if jwt.IsSaveCompanyIDOR(loginUser, tOther.CompanyID) {
		return errors.New(response.ErrorHandlerIDOR)
	}

	tx := conn.Begin()

	err = u.repository.Delete(tx, tOther)
	if err != nil {
		return errors.New(fmt.Sprint("failed to delete other: ", err))
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
