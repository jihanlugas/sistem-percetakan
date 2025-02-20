package phase

import (
	"github.com/jihanlugas/sistem-percetakan/jwt"
	"github.com/jihanlugas/sistem-percetakan/request"
	"github.com/jihanlugas/sistem-percetakan/response"
	"github.com/jihanlugas/sistem-percetakan/utils"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Handler struct {
	usecase Usecase
}

func NewHandler(usecase Usecase) Handler {
	return Handler{
		usecase: usecase,
	}
}

// Page
// @Tags Phase
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param req query request.PagePhase false "url query string"
// @Success      200  {object}	response.Response
// @Failure      500  {object}  response.Response
// @Router /phase [get]
func (h Handler) Page(c echo.Context) error {
	var err error

	loginUser, err := jwt.GetUserLoginInfo(c)
	if err != nil {
		return response.Error(http.StatusBadRequest, response.ErrorHandlerGetUserInfo, err, nil).SendJSON(c)
	}

	req := new(request.PagePhase)
	if err = c.Bind(req); err != nil {
		return response.Error(http.StatusBadRequest, response.ErrorHandlerBind, err, nil).SendJSON(c)
	}

	utils.TrimWhitespace(req)

	err = c.Validate(req)
	if err != nil {
		return response.Error(http.StatusBadRequest, response.ErrorHandlerFailedValidation, err, response.ValidationError(err)).SendJSON(c)
	}

	if req.CompanyID == "" {
		req.CompanyID = loginUser.CompanyID
	} else {
		if jwt.IsSaveCompanyIDOR(loginUser, req.CompanyID) {
			return response.Error(http.StatusBadRequest, response.ErrorHandlerIDOR, err, nil).SendJSON(c)
		}
	}

	data, count, err := h.usecase.Page(loginUser, *req)
	if err != nil {
		return response.Error(http.StatusBadRequest, err.Error(), err, nil).SendJSON(c)
	}

	return response.Success(http.StatusOK, response.SuccessHandler, response.PayloadPagination(req, data, count)).SendJSON(c)
}
