package photo

import (
	"errors"
	"fmt"
	"github.com/jihanlugas/sistem-percetakan/db"
	"github.com/jihanlugas/sistem-percetakan/model"
)

type Usecase interface {
	GetById(id string) (tPhoto model.Photo, err error)
	Upload() (tPhoto model.Photo, err error)
}

type usecase struct {
	repository Repository
}

func (u usecase) GetById(id string) (tPhoto model.Photo, err error) {
	conn, closeConn := db.GetConnection()
	defer closeConn()

	tPhoto, err = u.repository.GetById(conn, id)
	if err != nil {
		return tPhoto, errors.New(fmt.Sprint("failed to get order: ", err))
	}

	return tPhoto, err
}

func (u usecase) Upload() (tPhoto model.Photo, err error) {

	return tPhoto, err
}

func NewUsecase(repository Repository) Usecase {
	return &usecase{
		repository: repository,
	}
}
