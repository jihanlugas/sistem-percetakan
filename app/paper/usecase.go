package paper

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
	Page(loginUser jwt.UserLogin, req request.PagePaper) (vPapers []model.PaperView, count int64, err error)
	GetById(loginUser jwt.UserLogin, id string, preloads ...string) (vPaper model.PaperView, err error)
	Create(loginUser jwt.UserLogin, req request.CreatePaper) error
	Update(loginUser jwt.UserLogin, id string, req request.UpdatePaper) error
	Delete(loginUser jwt.UserLogin, id string) error
}

type usecase struct {
	repository Repository
}

func (u usecase) Page(loginUser jwt.UserLogin, req request.PagePaper) (vPapers []model.PaperView, count int64, err error) {
	conn, closeConn := db.GetConnection()
	defer closeConn()

	vPapers, count, err = u.repository.Page(conn, req)
	if err != nil {
		return vPapers, count, err
	}

	return vPapers, count, err
}

func (u usecase) GetById(loginUser jwt.UserLogin, id string, preloads ...string) (vPaper model.PaperView, err error) {
	conn, closeConn := db.GetConnection()
	defer closeConn()

	vPaper, err = u.repository.GetViewById(conn, id, preloads...)
	if err != nil {
		return vPaper, errors.New(fmt.Sprint("failed to get order: ", err))
	}

	if auth.IsSaveIDOR(loginUser, vPaper.CompanyID) {
		return vPaper, errors.New(response.ErrorHandlerIDOR)
	}

	return vPaper, err
}

func (u usecase) Create(loginUser jwt.UserLogin, req request.CreatePaper) error {
	var err error
	var tPaper model.Paper

	conn, closeConn := db.GetConnection()
	defer closeConn()

	if auth.IsSaveIDOR(loginUser, req.CompanyID) {
		return errors.New(response.ErrorHandlerIDOR)
	}

	tx := conn.Begin()

	tPaper = model.Paper{
		ID:                 utils.GetUniqueID(),
		CompanyID:          req.CompanyID,
		Name:               req.Name,
		Description:        req.Description,
		DefaultPrice:       req.DefaultPrice,
		DefaultPriceDuplex: req.DefaultPriceDuplex,
		CreateBy:           loginUser.UserID,
		UpdateBy:           loginUser.UserID,
	}

	err = u.repository.Create(tx, tPaper)
	if err != nil {
		return errors.New(fmt.Sprint("failed to create design: ", err))
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return err
}

func (u usecase) Update(loginUser jwt.UserLogin, id string, req request.UpdatePaper) error {
	var err error
	var tPaper model.Paper

	conn, closeConn := db.GetConnection()
	defer closeConn()

	tPaper, err = u.repository.GetTableById(conn, id)
	if err != nil {
		return errors.New(fmt.Sprint("failed to get design: ", err))
	}

	if auth.IsSaveIDOR(loginUser, tPaper.CompanyID) {
		return errors.New(response.ErrorHandlerIDOR)
	}

	tx := conn.Begin()

	tPaper.Name = req.Name
	tPaper.Description = req.Description
	tPaper.DefaultPrice = req.DefaultPrice
	tPaper.DefaultPriceDuplex = req.DefaultPriceDuplex
	tPaper.UpdateBy = loginUser.UserID
	err = u.repository.Save(tx, tPaper)
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
	var tPaper model.Paper

	conn, closeConn := db.GetConnection()
	defer closeConn()

	tPaper, err = u.repository.GetTableById(conn, id)
	if err != nil {
		return errors.New(fmt.Sprint("failed to get design: ", err))
	}

	if auth.IsSaveIDOR(loginUser, tPaper.CompanyID) {
		return errors.New(response.ErrorHandlerIDOR)
	}

	tx := conn.Begin()

	err = u.repository.Delete(tx, tPaper)
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
