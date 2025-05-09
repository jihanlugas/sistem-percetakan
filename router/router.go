package router

import (
	"encoding/json"
	"fmt"
	"github.com/jihanlugas/sistem-percetakan/app/auth"
	"github.com/jihanlugas/sistem-percetakan/app/company"
	"github.com/jihanlugas/sistem-percetakan/app/customer"
	"github.com/jihanlugas/sistem-percetakan/app/dashboard"
	"github.com/jihanlugas/sistem-percetakan/app/finishing"
	"github.com/jihanlugas/sistem-percetakan/app/order"
	"github.com/jihanlugas/sistem-percetakan/app/orderphase"
	"github.com/jihanlugas/sistem-percetakan/app/paper"
	"github.com/jihanlugas/sistem-percetakan/app/phase"
	"github.com/jihanlugas/sistem-percetakan/app/photo"
	"github.com/jihanlugas/sistem-percetakan/app/print"
	"github.com/jihanlugas/sistem-percetakan/app/transaction"
	"github.com/jihanlugas/sistem-percetakan/app/user"
	"github.com/jihanlugas/sistem-percetakan/app/usercompany"
	"github.com/jihanlugas/sistem-percetakan/config"
	"github.com/jihanlugas/sistem-percetakan/constant"
	"github.com/jihanlugas/sistem-percetakan/db"
	"github.com/jihanlugas/sistem-percetakan/jwt"
	"github.com/jihanlugas/sistem-percetakan/model"
	"github.com/jihanlugas/sistem-percetakan/response"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"net/http"
)

func Init() *echo.Echo {
	router := websiteRouter()

	photoRepository := photo.NewRepository()
	userRepository := user.NewRepository()
	companyRepository := company.NewRepository()
	usercompanyRepository := usercompany.NewRepository()
	orderRepository := order.NewRepository()
	printRepository := print.NewRepository()
	finishingRepository := finishing.NewRepository()
	orderphaseRepository := orderphase.NewRepository()
	customerRepository := customer.NewRepository()
	paperRepository := paper.NewRepository()
	phaseRepository := phase.NewRepository()
	transactionRepository := transaction.NewRepository()

	authUsecase := auth.NewUsecase(userRepository, companyRepository, usercompanyRepository)
	photoUsecase := photo.NewUsecase(photoRepository)
	userUsecase := user.NewUsecase(userRepository, usercompanyRepository)
	orderUsecase := order.NewUsecase(orderRepository, printRepository, finishingRepository, orderphaseRepository, customerRepository, phaseRepository, transactionRepository)
	customerUsecase := customer.NewUsecase(customerRepository)
	paperUsecase := paper.NewUsecase(paperRepository)
	phaseUsecase := phase.NewUsecase(phaseRepository)
	printUsecase := print.NewUsecase(printRepository)
	finishingUsecase := finishing.NewUsecase(finishingRepository)
	transactionUsecase := transaction.NewUsecase(transactionRepository)
	dashboardUsecase := dashboard.NewUsecase(orderRepository, transactionRepository)
	companyUsecase := company.NewUsecase(companyRepository, usercompanyRepository)

	authHandler := auth.NewHandler(authUsecase)
	photoHandler := photo.NewHandler(photoUsecase)
	companyHandler := company.NewHandler(companyUsecase)
	userHandler := user.NewHandler(userUsecase)
	orderHandler := order.NewHandler(orderUsecase)
	customerHandler := customer.NewHandler(customerUsecase)
	paperHandler := paper.NewHandler(paperUsecase)
	phaseHandler := phase.NewHandler(phaseUsecase)
	printHandler := print.NewHandler(printUsecase)
	finishingHandler := finishing.NewHandler(finishingUsecase)
	transactionHandler := transaction.NewHandler(transactionUsecase)
	dashboardHandler := dashboard.NewHandler(dashboardUsecase)

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

	routerPhoto := router.Group("/photo")
	routerPhoto.GET("/:id", photoHandler.GetById)

	routerDashboard := router.Group("/dashboard", checkTokenMiddleware)
	routerDashboard.GET("/:id", dashboardHandler.GetDashboardById)

	routerUser := router.Group("/user", checkTokenMiddleware)
	routerUser.GET("", userHandler.Page)
	routerUser.POST("", userHandler.Create)
	routerUser.POST("/change-password", userHandler.ChangePassword)
	routerUser.PUT("/:id", userHandler.Update)
	routerUser.GET("/:id", userHandler.GetById)
	routerUser.DELETE("/:id", userHandler.Delete)

	routerCompany := router.Group("/company", checkTokenMiddleware)
	routerCompany.PUT("/:id", companyHandler.Update)

	routerOrder := router.Group("/order", checkTokenMiddleware)
	routerOrder.GET("", orderHandler.Page)
	routerOrder.POST("", orderHandler.Create)
	routerOrder.GET("/:id", orderHandler.GetById)
	routerOrder.GET("/:id/spk", orderHandler.GenerateSpk)
	routerOrder.GET("/:id/invoice", orderHandler.GenerateInvoice)
	routerOrder.PUT("/:id", orderHandler.Update)
	routerOrder.POST("/:id/add-phase", orderHandler.AddPhase)
	routerOrder.POST("/:id/add-transaction", orderHandler.AddTransaction)
	routerOrder.DELETE("/:id", orderHandler.Delete)

	//router.GET("/order/:id/spk", orderHandler.GenerateSpk)
	//router.GET("/order/:id/invoice", orderHandler.GenerateInvoice)

	routerCustomer := router.Group("/customer", checkTokenMiddleware)
	routerCustomer.GET("", customerHandler.Page)
	routerCustomer.POST("", customerHandler.Create)
	routerCustomer.PUT("/:id", customerHandler.Update)
	routerCustomer.GET("/:id", customerHandler.GetById)
	routerCustomer.DELETE("/:id", customerHandler.Delete)

	routerPaper := router.Group("/paper", checkTokenMiddleware)
	routerPaper.GET("", paperHandler.Page)
	routerPaper.POST("", paperHandler.Create)
	routerPaper.PUT("/:id", paperHandler.Update)
	routerPaper.GET("/:id", paperHandler.GetById)
	routerPaper.DELETE("/:id", paperHandler.Delete)

	routerPhase := router.Group("/phase", checkTokenMiddleware)
	routerPhase.GET("", phaseHandler.Page)

	routerPrint := router.Group("/print", checkTokenMiddleware)
	routerPrint.GET("", printHandler.Page)
	routerPrint.POST("", printHandler.Create)
	routerPrint.PUT("/:id", printHandler.Update)
	routerPrint.GET("/:id", printHandler.GetById)
	routerPrint.DELETE("/:id", printHandler.Delete)
	routerPrint.GET("/:id/spk", printHandler.GenerateSpk)

	routerFinishing := router.Group("/finishing", checkTokenMiddleware)
	routerFinishing.GET("", finishingHandler.Page)
	routerFinishing.POST("", finishingHandler.Create)
	routerFinishing.PUT("/:id", finishingHandler.Update)
	routerFinishing.GET("/:id", finishingHandler.GetById)
	routerFinishing.DELETE("/:id", finishingHandler.Delete)

	routerTransaction := router.Group("/transaction", checkTokenMiddleware)
	routerTransaction.GET("", transactionHandler.Page)
	routerTransaction.POST("", transactionHandler.Create)
	routerTransaction.PUT("/:id", transactionHandler.Update)
	routerTransaction.GET("/:id", transactionHandler.GetById)
	routerTransaction.DELETE("/:id", transactionHandler.Delete)

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
