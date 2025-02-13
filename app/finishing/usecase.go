package finishing

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
	Page(loginUser jwt.UserLogin, req request.PageFinishing) (vFinishings []model.FinishingView, count int64, err error)
	GetById(loginUser jwt.UserLogin, id string, preloads ...string) (vFinishing model.FinishingView, err error)
	Create(loginUser jwt.UserLogin, req request.CreateFinishing) error
	Update(loginUser jwt.UserLogin, id string, req request.UpdateFinishing) error
	Delete(loginUser jwt.UserLogin, id string) error
}

type usecase struct {
	repository Repository
}

func (u usecase) Page(loginUser jwt.UserLogin, req request.PageFinishing) (vFinishings []model.FinishingView, count int64, err error) {
	conn, closeConn := db.GetConnection()
	defer closeConn()

	vFinishings, count, err = u.repository.Page(conn, req)
	if err != nil {
		return vFinishings, count, err
	}

	return vFinishings, count, err
}

func (u usecase) GetById(loginUser jwt.UserLogin, id string, preloads ...string) (vFinishing model.FinishingView, err error) {
	conn, closeConn := db.GetConnection()
	defer closeConn()

	vFinishing, err = u.repository.GetViewById(conn, id, preloads...)
	if err != nil {
		return vFinishing, errors.New(fmt.Sprint("failed to get order: ", err))
	}

	if auth.IsSaveIDOR(loginUser, vFinishing.CompanyID) {
		return vFinishing, errors.New(response.ErrorHandlerIDOR)
	}

	return vFinishing, err
}

func (u usecase) Create(loginUser jwt.UserLogin, req request.CreateFinishing) error {
	var err error
	var tFinishing model.Finishing

	conn, closeConn := db.GetConnection()
	defer closeConn()

	if auth.IsSaveIDOR(loginUser, req.CompanyID) {
		return errors.New(response.ErrorHandlerIDOR)
	}

	tx := conn.Begin()

	tFinishing = model.Finishing{
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

	err = u.repository.Create(tx, tFinishing)
	if err != nil {
		return errors.New(fmt.Sprint("failed to create finishing: ", err))
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return err
}

func (u usecase) Update(loginUser jwt.UserLogin, id string, req request.UpdateFinishing) error {
	var err error
	var tFinishing model.Finishing

	conn, closeConn := db.GetConnection()
	defer closeConn()

	tFinishing, err = u.repository.GetTableById(conn, id)
	if err != nil {
		return errors.New(fmt.Sprint("failed to get finishing: ", err))
	}

	if auth.IsSaveIDOR(loginUser, tFinishing.CompanyID) {
		return errors.New(response.ErrorHandlerIDOR)
	}

	tx := conn.Begin()

	tFinishing.Name = req.Name
	tFinishing.Description = req.Description
	tFinishing.Qty = req.Qty
	tFinishing.Price = req.Price
	tFinishing.Total = req.Total
	tFinishing.UpdateBy = loginUser.UserID
	err = u.repository.Save(tx, tFinishing)
	if err != nil {
		return errors.New(fmt.Sprint("failed to update finishing: ", err))
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return err
}

func (u usecase) Delete(loginUser jwt.UserLogin, id string) error {
	var err error
	var tFinishing model.Finishing

	conn, closeConn := db.GetConnection()
	defer closeConn()

	tFinishing, err = u.repository.GetTableById(conn, id)
	if err != nil {
		return errors.New(fmt.Sprint("failed to get finishing: ", err))
	}

	if auth.IsSaveIDOR(loginUser, tFinishing.CompanyID) {
		return errors.New(response.ErrorHandlerIDOR)
	}

	tx := conn.Begin()

	err = u.repository.Delete(tx, tFinishing)
	if err != nil {
		return errors.New(fmt.Sprint("failed to delete finishing: ", err))
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
