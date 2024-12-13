package photo

import "github.com/jihanlugas/sistem-percetakan/model"

type Usecase interface {
	Upload() (tPhoto model.Photo, err error)
}

type usecase struct {
	repo Repository
}

func (u usecase) Upload() (tPhoto model.Photo, err error) {

	return tPhoto, err
}

func NewUsecase(repo Repository) Usecase {
	return &usecase{
		repo: repo,
	}
}
