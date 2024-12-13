package auth

import (
	"errors"
	"github.com/jihanlugas/sistem-percetakan/app/company"
	"github.com/jihanlugas/sistem-percetakan/app/user"
	"github.com/jihanlugas/sistem-percetakan/app/usercompany"
	"github.com/jihanlugas/sistem-percetakan/config"
	"github.com/jihanlugas/sistem-percetakan/constant"
	"github.com/jihanlugas/sistem-percetakan/cryption"
	"github.com/jihanlugas/sistem-percetakan/db"
	"github.com/jihanlugas/sistem-percetakan/jwt"
	"github.com/jihanlugas/sistem-percetakan/model"
	"github.com/jihanlugas/sistem-percetakan/request"
	"github.com/jihanlugas/sistem-percetakan/utils"
	"time"
)

type Usecase interface {
	SignIn(req request.Signin) (token string, userLogin jwt.UserLogin, err error)
	RefreshToken(userLogin jwt.UserLogin) (token string, err error)
	Init(userLogin jwt.UserLogin) (vUser model.UserView, vCompany model.CompanyView, err error)
}

type usecase struct {
	userRepository        user.Repository
	companyRepository     company.Repository
	usercompanyRepository usercompany.Repository
}

func (u usecase) SignIn(req request.Signin) (token string, userLogin jwt.UserLogin, err error) {

	var tUser model.User
	var tCompany model.Company
	var tUsercompany model.Usercompany

	conn, closeConn := db.GetConnection()
	defer closeConn()

	if utils.IsValidEmail(req.Username) {
		tUser, err = u.userRepository.GetByEmail(conn, req.Username)
	} else {
		tUser, err = u.userRepository.GetByUsername(conn, req.Username)
	}

	if err != nil {
		return "", userLogin, err
	}

	err = cryption.CheckAES64(req.Passwd, tUser.Passwd)
	if err != nil {
		return "", userLogin, errors.New("invalid username or password")
	}

	if !tUser.IsActive {
		return "", userLogin, errors.New("user not active")
	}

	if tUser.Role != constant.RoleAdmin {
		tUsercompany, err = u.usercompanyRepository.GetCompanyDefaultByUserId(conn, tUser.ID)
		if err != nil {
			return "", userLogin, errors.New("usercompany not found : " + err.Error())
		}

		tCompany, err = u.companyRepository.GetById(conn, tUsercompany.CompanyID)
		if err != nil {
			return "", userLogin, errors.New("company not found : " + err.Error())
		}
	}

	now := time.Now()
	tx := conn.Begin()

	tUser.LastLoginDt = &now
	tUser.UpdateBy = tUser.ID
	err = u.userRepository.Update(tx, model.User{
		ID:          tUser.ID,
		LastLoginDt: &now,
		UpdateBy:    tUser.ID,
	})
	if err != nil {
		return "", userLogin, err
	}

	err = tx.Commit().Error
	if err != nil {
		return "", userLogin, err
	}

	expiredAt := time.Now().Add(time.Minute * time.Duration(config.AuthTokenExpiredMinute))
	userLogin.ExpiredDt = expiredAt
	userLogin.UserID = tUser.ID
	userLogin.Role = tUser.Role
	userLogin.PassVersion = tUser.PassVersion
	userLogin.CompanyID = tCompany.ID
	userLogin.UsercompanyID = tUsercompany.ID
	token, err = jwt.CreateToken(userLogin)
	if err != nil {
		return "", userLogin, err
	}

	return token, userLogin, err
}

func (u usecase) RefreshToken(userLogin jwt.UserLogin) (token string, err error) {
	userLogin.ExpiredDt = time.Now().Add(time.Minute * time.Duration(config.AuthTokenExpiredMinute))

	token, err = jwt.CreateToken(userLogin)
	if err != nil {
		return token, err
	}

	return token, err
}

func (u usecase) Init(userLogin jwt.UserLogin) (vUser model.UserView, vCompany model.CompanyView, err error) {
	conn, closeConn := db.GetConnection()
	defer closeConn()

	vUser, err = u.userRepository.GetViewById(conn, userLogin.UserID)
	if err != nil {
		return vUser, vCompany, err
	}

	if userLogin.Role != constant.RoleAdmin {
		vCompany, err = u.companyRepository.GetViewById(conn, userLogin.CompanyID)
		if err != nil {
			return vUser, vCompany, err
		}
	}

	return vUser, vCompany, err
}

func NewUsecase(userRepository user.Repository, companyRepository company.Repository, usercompanyRepository usercompany.Repository) Usecase {
	return usecase{
		userRepository:        userRepository,
		companyRepository:     companyRepository,
		usercompanyRepository: usercompanyRepository,
	}
}

func IsSaveIDOR(loginUser jwt.UserLogin, companyId string) bool {
	if loginUser.Role != constant.RoleAdmin {
		if loginUser.CompanyID != companyId {
			return true
		}
	}

	return false
}
