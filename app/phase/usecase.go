package phase

import (
	"github.com/jihanlugas/sistem-percetakan/db"
	"github.com/jihanlugas/sistem-percetakan/jwt"
	"github.com/jihanlugas/sistem-percetakan/model"
	"github.com/jihanlugas/sistem-percetakan/request"
)

type Usecase interface {
	Page(loginUser jwt.UserLogin, req request.PagePhase) (vPhases []model.PhaseView, count int64, err error)
}

type usecase struct {
	repository Repository
}

func (u usecase) Page(loginUser jwt.UserLogin, req request.PagePhase) (vPhases []model.PhaseView, count int64, err error) {
	conn, closeConn := db.GetConnection()
	defer closeConn()

	vPhases, count, err = u.repository.Page(conn, req)
	if err != nil {
		return vPhases, count, err
	}

	return vPhases, count, err
}

func NewUsecase(repository Repository) Usecase {
	return &usecase{
		repository: repository,
	}
}
