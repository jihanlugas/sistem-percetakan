package auth

import (
	"github.com/jihanlugas/sistem-percetakan/jwt"
	"github.com/jihanlugas/sistem-percetakan/request"
	"github.com/jihanlugas/sistem-percetakan/response"
	"github.com/jihanlugas/sistem-percetakan/utils"
	"github.com/labstack/echo/v4"
	"net/http"
)

type handler struct {
	usecase Usecase
}

func NewHandler(usecase Usecase) handler {
	return handler{
		usecase: usecase,
	}
}

// SignIn
// @Tags Authentication
// @Accept json
// @Produce json
// @Param req body request.Signin true "json req body"
// @Success      200  {object}	response.Response
// @Failure      500  {object}  response.Response
// @Router /auth/sign-in [post]
func (h handler) SignIn(c echo.Context) error {
	var err error

	req := new(request.Signin)
	if err = c.Bind(req); err != nil {
		return response.Error(http.StatusBadRequest, response.ErrorHandlerBind, err, nil).SendJSON(c)
	}

	utils.TrimWhitespace(req)

	err = c.Validate(req)
	if err != nil {
		return response.Error(http.StatusBadRequest, response.ErrorHandlerFailedValidation, err, response.ValidationError(err)).SendJSON(c)
	}

	token, userLogin, err := h.usecase.SignIn(*req)
	if err != nil {
		return response.Error(http.StatusBadRequest, err.Error(), err, nil).SendJSON(c)
	}

	return response.Success(http.StatusOK, response.SuccessHandler, response.Payload{
		"token":     token,
		"userLogin": userLogin,
	}).SendJSON(c)
}

// SignOut Sign out user
// @Tags Authentication
// @Accept json
// @Produce json
// // @Param req body request.Signin true "json req body"
// @Success      200  {object}	response.Response
// @Failure      500  {object}  response.Response
// @Router /auth/sign-out [get]
func (h handler) SignOut(c echo.Context) error {
	return response.Success(http.StatusOK, response.SuccessHandler, nil).SendJSON(c)
}

// RefreshToken
// @Tags Authentication
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success      200  {object}	response.Response
// @Failure      500  {object}  response.Response
// @Router /auth/refresh-token [get]
func (h handler) RefreshToken(c echo.Context) error {
	var err error

	loginUser, err := jwt.GetUserLoginInfo(c)
	if err != nil {
		return response.Error(http.StatusBadRequest, response.ErrorHandlerBind, err, nil).SendJSON(c)
	}

	token, err := h.usecase.RefreshToken(loginUser)
	if err != nil {
		return response.Error(http.StatusBadRequest, err.Error(), err, nil).SendJSON(c)
	}

	return response.Success(http.StatusOK, response.SuccessHandler, response.Payload{
		"token": token,
	}).SendJSON(c)
}

// Init
// @Tags Authentication
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success      200  {object}	response.Response
// @Failure      500  {object}  response.Response
// @Router /auth/init [get]
func (h handler) Init(c echo.Context) error {
	var err error

	loginUser, err := jwt.GetUserLoginInfo(c)
	if err != nil {
		return response.Error(http.StatusBadRequest, response.ErrorHandlerGetUserInfo, err, nil).SendJSON(c)
	}

	user, company, err := h.usecase.Init(loginUser)
	if err != nil {
		return response.Error(http.StatusBadRequest, err.Error(), err, nil).SendJSON(c)
	}

	resUser := response.User(user)
	resCompany := response.Company(company)

	res := response.Init{
		User:    &resUser,
		Company: &resCompany,
	}

	return response.Success(http.StatusOK, response.SuccessHandler, res).SendJSON(c)
}
