package print

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
	Page(loginUser jwt.UserLogin, req request.PagePrint) (vPrints []model.PrintView, count int64, err error)
	GetById(loginUser jwt.UserLogin, id string, preloads ...string) (vPrint model.PrintView, err error)
	Create(loginUser jwt.UserLogin, req request.CreatePrint) error
	Update(loginUser jwt.UserLogin, id string, req request.UpdatePrint) error
	Delete(loginUser jwt.UserLogin, id string) error
}

type usecase struct {
	repository Repository
}

func (u usecase) Page(loginUser jwt.UserLogin, req request.PagePrint) (vPrints []model.PrintView, count int64, err error) {
	conn, closeConn := db.GetConnection()
	defer closeConn()

	vPrints, count, err = u.repository.Page(conn, req)
	if err != nil {
		return vPrints, count, err
	}

	return vPrints, count, err
}

func (u usecase) GetById(loginUser jwt.UserLogin, id string, preloads ...string) (vPrint model.PrintView, err error) {
	conn, closeConn := db.GetConnection()
	defer closeConn()

	vPrint, err = u.repository.GetViewById(conn, id, preloads...)
	if err != nil {
		return vPrint, errors.New(fmt.Sprint("failed to get order: ", err))
	}

	if jwt.IsSaveCompanyIDOR(loginUser, vPrint.CompanyID) {
		return vPrint, errors.New(response.ErrorHandlerIDOR)
	}

	return vPrint, err
}

func (u usecase) Create(loginUser jwt.UserLogin, req request.CreatePrint) error {
	var err error
	var tPrint model.Print

	conn, closeConn := db.GetConnection()
	defer closeConn()

	if jwt.IsSaveCompanyIDOR(loginUser, req.CompanyID) {
		return errors.New(response.ErrorHandlerIDOR)
	}

	tx := conn.Begin()

	tPrint = model.Print{
		ID:          utils.GetUniqueID(),
		CompanyID:   req.CompanyID,
		OrderID:     req.OrderID,
		PaperID:     req.PaperID,
		Name:        req.Name,
		Description: req.Description,
		IsDuplex:    req.IsDuplex,
		PageCount:   req.PageCount,
		Qty:         req.Qty,
		Price:       req.Price,
		Total:       req.Total,
		CreateBy:    loginUser.UserID,
		UpdateBy:    loginUser.UserID,
	}

	err = u.repository.Create(tx, tPrint)
	if err != nil {
		return errors.New(fmt.Sprint("failed to create print: ", err))
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return err
}

func (u usecase) Update(loginUser jwt.UserLogin, id string, req request.UpdatePrint) error {
	var err error
	var tPrint model.Print

	conn, closeConn := db.GetConnection()
	defer closeConn()

	tPrint, err = u.repository.GetTableById(conn, id)
	if err != nil {
		return errors.New(fmt.Sprint("failed to get print: ", err))
	}

	if jwt.IsSaveCompanyIDOR(loginUser, tPrint.CompanyID) {
		return errors.New(response.ErrorHandlerIDOR)
	}

	tx := conn.Begin()

	tPrint.PaperID = req.PaperID
	tPrint.Name = req.Name
	tPrint.Description = req.Description
	tPrint.IsDuplex = req.IsDuplex
	tPrint.PageCount = req.PageCount
	tPrint.Qty = req.Qty
	tPrint.Price = req.Price
	tPrint.Total = req.Total
	tPrint.UpdateBy = loginUser.UserID
	err = u.repository.Save(tx, tPrint)
	if err != nil {
		return errors.New(fmt.Sprint("failed to update print: ", err))
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return err
}

func (u usecase) Delete(loginUser jwt.UserLogin, id string) error {
	var err error
	var tPrint model.Print

	conn, closeConn := db.GetConnection()
	defer closeConn()

	tPrint, err = u.repository.GetTableById(conn, id)
	if err != nil {
		return errors.New(fmt.Sprint("failed to get print: ", err))
	}

	if jwt.IsSaveCompanyIDOR(loginUser, tPrint.CompanyID) {
		return errors.New(response.ErrorHandlerIDOR)
	}

	tx := conn.Begin()

	err = u.repository.Delete(tx, tPrint)
	if err != nil {
		return errors.New(fmt.Sprint("failed to delete print: ", err))
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
