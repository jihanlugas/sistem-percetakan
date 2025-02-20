package design

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
	Page(loginUser jwt.UserLogin, req request.PageDesign) (vDesigns []model.DesignView, count int64, err error)
	GetById(loginUser jwt.UserLogin, id string, preloads ...string) (vDesign model.DesignView, err error)
	Create(loginUser jwt.UserLogin, req request.CreateDesign) error
	Update(loginUser jwt.UserLogin, id string, req request.UpdateDesign) error
	Delete(loginUser jwt.UserLogin, id string) error
}

type usecase struct {
	repository Repository
}

func (u usecase) Page(loginUser jwt.UserLogin, req request.PageDesign) (vDesigns []model.DesignView, count int64, err error) {
	conn, closeConn := db.GetConnection()
	defer closeConn()

	vDesigns, count, err = u.repository.Page(conn, req)
	if err != nil {
		return vDesigns, count, err
	}

	return vDesigns, count, err
}

func (u usecase) GetById(loginUser jwt.UserLogin, id string, preloads ...string) (vDesign model.DesignView, err error) {
	conn, closeConn := db.GetConnection()
	defer closeConn()

	vDesign, err = u.repository.GetViewById(conn, id, preloads...)
	if err != nil {
		return vDesign, errors.New(fmt.Sprint("failed to get order: ", err))
	}

	if jwt.IsSaveCompanyIDOR(loginUser, vDesign.CompanyID) {
		return vDesign, errors.New(response.ErrorHandlerIDOR)
	}

	return vDesign, err
}

func (u usecase) Create(loginUser jwt.UserLogin, req request.CreateDesign) error {
	var err error
	var tDesign model.Design

	conn, closeConn := db.GetConnection()
	defer closeConn()

	if jwt.IsSaveCompanyIDOR(loginUser, req.CompanyID) {
		return errors.New(response.ErrorHandlerIDOR)
	}

	tx := conn.Begin()

	tDesign = model.Design{
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

	err = u.repository.Create(tx, tDesign)
	if err != nil {
		return errors.New(fmt.Sprint("failed to create design: ", err))
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return err
}

func (u usecase) Update(loginUser jwt.UserLogin, id string, req request.UpdateDesign) error {
	var err error
	var tDesign model.Design

	conn, closeConn := db.GetConnection()
	defer closeConn()

	tDesign, err = u.repository.GetTableById(conn, id)
	if err != nil {
		return errors.New(fmt.Sprint("failed to get design: ", err))
	}

	if jwt.IsSaveCompanyIDOR(loginUser, tDesign.CompanyID) {
		return errors.New(response.ErrorHandlerIDOR)
	}

	tx := conn.Begin()

	tDesign.Name = req.Name
	tDesign.Description = req.Description
	tDesign.Qty = req.Qty
	tDesign.Price = req.Price
	tDesign.Total = req.Total
	tDesign.UpdateBy = loginUser.UserID
	err = u.repository.Save(tx, tDesign)
	if err != nil {
		return errors.New(fmt.Sprint("failed to update design: ", err))
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return err
}

func (u usecase) Delete(loginUser jwt.UserLogin, id string) error {
	var err error
	var tDesign model.Design

	conn, closeConn := db.GetConnection()
	defer closeConn()

	tDesign, err = u.repository.GetTableById(conn, id)
	if err != nil {
		return errors.New(fmt.Sprint("failed to get design: ", err))
	}

	if jwt.IsSaveCompanyIDOR(loginUser, tDesign.CompanyID) {
		return errors.New(response.ErrorHandlerIDOR)
	}

	tx := conn.Begin()

	err = u.repository.Delete(tx, tDesign)
	if err != nil {
		return errors.New(fmt.Sprint("failed to delete design: ", err))
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
