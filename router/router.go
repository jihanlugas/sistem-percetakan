package router

import (
	"encoding/json"
	"fmt"
	"github.com/jihanlugas/sistem-percetakan/app/auth"
	"github.com/jihanlugas/sistem-percetakan/app/company"
	"github.com/jihanlugas/sistem-percetakan/app/user"
	"github.com/jihanlugas/sistem-percetakan/app/usercompany"
	"github.com/jihanlugas/sistem-percetakan/config"
	"github.com/jihanlugas/sistem-percetakan/constant"
	"github.com/jihanlugas/sistem-percetakan/db"
	"github.com/jihanlugas/sistem-percetakan/jwt"
	"github.com/jihanlugas/sistem-percetakan/model"
	"github.com/jihanlugas/sistem-percetakan/response"
	"github.com/labstack/echo/v4"
	"net/http"

	echoSwagger "github.com/swaggo/echo-swagger"
)

func Init() *echo.Echo {
	router := websiteRouter()

	userRepository := user.NewRepository()
	companyRepository := company.NewRepository()
	usercompanyRepository := usercompany.NewRepository()

	authUsecase := auth.NewUsecase(userRepository, companyRepository, usercompanyRepository)

	authHandler := auth.NewHandler(authUsecase)

	if config.Debug {
		router.GET("/", func(c echo.Context) error {
			return response.Success(http.StatusOK, "Welcome", nil).SendJSON(c)
		})
		router.GET("/swg/*", echoSwagger.WrapHandler)
	}

	router.Static("/storage", "storage")

	routerAuth := router.Group("/auth")
	routerAuth.POST("/sign-in", authHandler.SignIn)
	routerAuth.POST("/sign-out", authHandler.SignOut)
	routerAuth.GET("/init", authHandler.Init, checkTokenMiddleware)
	routerAuth.GET("/refresh-token", authHandler.RefreshToken, checkTokenMiddleware)

	return router
}

func httpErrorHandler(err error, c echo.Context) {
	var errorResponse *response.Response
	code := http.StatusInternalServerError
	switch e := err.(type) {
	case *echo.HTTPError:
		// Handle pada saat URL yang di request tidak ada. atau ada kesalahan server.
		code = e.Code
		errorResponse = &response.Response{
			Status:  false,
			Message: fmt.Sprintf("%v", e.Message),
			Code:    code,
		}
	case *response.Response:
		errorResponse = e
	default:
		// Handle error dari panic
		code = http.StatusInternalServerError
		if config.Debug {
			errorResponse = &response.Response{
				Status:  false,
				Message: err.Error(),
				Code:    http.StatusInternalServerError,
			}
		} else {
			errorResponse = &response.Response{
				Status:  false,
				Message: response.ErrorInternalServer,
				Code:    http.StatusInternalServerError,
			}
		}
	}

	js, err := json.Marshal(errorResponse)
	if err == nil {
		_ = c.Blob(code, echo.MIMEApplicationJSONCharsetUTF8, js)
	} else {
		b := []byte("{status: false, code: 500, message: \"unresolved error\"}")
		_ = c.Blob(code, echo.MIMEApplicationJSONCharsetUTF8, b)
	}
}

func checkTokenMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var err error

		userLogin, err := jwt.ExtractClaims(c.Request().Header.Get(constant.AuthHeaderKey))
		if err != nil {
			return response.ErrorForce(http.StatusUnauthorized, err.Error()).SendJSON(c)
		}

		conn, closeConn := db.GetConnection()
		defer closeConn()

		var user model.User
		err = conn.Where("id = ? ", userLogin.UserID).First(&user).Error
		if err != nil {
			return response.ErrorForce(http.StatusUnauthorized, response.ErrorMiddlewareUserNotFound).SendJSON(c)
		}

		if user.PassVersion != userLogin.PassVersion {
			return response.ErrorForce(http.StatusUnauthorized, response.ErrorMiddlewarePassVersion).SendJSON(c)
		}

		c.Set(constant.TokenUserContext, userLogin)
		return next(c)
	}
}
