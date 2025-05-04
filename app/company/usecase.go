package company

import (
	"errors"
	"fmt"
	"github.com/jihanlugas/sistem-percetakan/app/usercompany"
	"github.com/jihanlugas/sistem-percetakan/constant"
	"github.com/jihanlugas/sistem-percetakan/db"
	"github.com/jihanlugas/sistem-percetakan/jwt"
	"github.com/jihanlugas/sistem-percetakan/model"
	"github.com/jihanlugas/sistem-percetakan/request"
	"github.com/jihanlugas/sistem-percetakan/response"
	"github.com/jihanlugas/sistem-percetakan/utils"
)

type Usecase interface {
	Update(loginUser jwt.UserLogin, id string, req request.UpdateCompany) error
}

type usecase struct {
	companyRepository     Repository
	usercompanyRepository usercompany.Repository
}

func NewUsecase(companyRepository Repository, usercompanyRepository usercompany.Repository) Usecase {
	return &usecase{
		companyRepository:     companyRepository,
		usercompanyRepository: usercompanyRepository,
	}
}

func (u usecase) Update(loginUser jwt.UserLogin, id string, req request.UpdateCompany) error {
	var err error
	var tCompany model.Company

	conn, closeConn := db.GetConnection()
	defer closeConn()

	tCompany, err = u.companyRepository.GetTableById(conn, id)
	if err != nil {
		return errors.New(fmt.Sprint("failed to get company: ", err))
	}

	switch loginUser.Role {
	case constant.RoleAdmin:
	case constant.RoleUser:
		return errors.New("role not allowed")
	case constant.RoleUseradmin:
		vUsercompany, err := u.usercompanyRepository.GetViewByUserIdAndCompanyId(conn, loginUser.UserID, loginUser.CompanyID)
		if err != nil {
			return errors.New(fmt.Sprint("failed to get user company: ", err))
		}

		if jwt.IsSaveCompanyIDOR(loginUser, vUsercompany.CompanyID) {
			return errors.New(response.ErrorHandlerIDOR)
		}
	default:
		return errors.New("not allowed")
	}

	tx := conn.Begin()

	tCompany.Name = req.Name
	tCompany.Description = req.Description
	tCompany.Email = req.Email
	tCompany.Address = req.Address
	tCompany.InvoiceNote = req.InvoiceNote
	tCompany.PhoneNumber = utils.FormatPhoneTo62(req.PhoneNumber)
	tCompany.UpdateBy = loginUser.UserID
	err = u.companyRepository.Save(tx, tCompany)
	if err != nil {
		return errors.New(fmt.Sprint("failed to update company: ", err))
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return err

}
