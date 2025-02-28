package router

import (
	"encoding/json"
	"fmt"
	"github.com/jihanlugas/sistem-percetakan/app/auth"
	"github.com/jihanlugas/sistem-percetakan/app/company"
	"github.com/jihanlugas/sistem-percetakan/app/customer"
	"github.com/jihanlugas/sistem-percetakan/app/design"
	"github.com/jihanlugas/sistem-percetakan/app/finishing"
	"github.com/jihanlugas/sistem-percetakan/app/order"
	"github.com/jihanlugas/sistem-percetakan/app/orderphase"
	"github.com/jihanlugas/sistem-percetakan/app/other"
	"github.com/jihanlugas/sistem-percetakan/app/paper"
	"github.com/jihanlugas/sistem-percetakan/app/phase"
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

	userRepository := user.NewRepository()
	companyRepository := company.NewRepository()
	usercompanyRepository := usercompany.NewRepository()
	orderRepository := order.NewRepository()
	designRepository := design.NewRepository()
	printRepository := print.NewRepository()
	finishingRepository := finishing.NewRepository()
	otherRepository := other.NewRepository()
	orderphaseRepository := orderphase.NewRepository()
	customerRepository := customer.NewRepository()
	paperRepository := paper.NewRepository()
	phaseRepository := phase.NewRepository()
	transactionRepository := transaction.NewRepository()

	authUsecase := auth.NewUsecase(userRepository, companyRepository, usercompanyRepository)
	userUsecase := user.NewUsecase(userRepository, usercompanyRepository)
	orderUsecase := order.NewUsecase(orderRepository, designRepository, printRepository, finishingRepository, otherRepository, orderphaseRepository, customerRepository, phaseRepository, transactionRepository)
	customerUsecase := customer.NewUsecase(customerRepository)
	paperUsecase := paper.NewUsecase(paperRepository)
	phaseUsecase := phase.NewUsecase(phaseRepository)
	designUsecase := design.NewUsecase(designRepository)
	printUsecase := print.NewUsecase(printRepository)
	finishingUsecase := finishing.NewUsecase(finishingRepository)
	otherUsecase := other.NewUsecase(otherRepository)
	transactionUsecase := transaction.NewUsecase(transactionRepository)

	authHandler := auth.NewHandler(authUsecase)
	userHandler := user.NewHandler(userUsecase)
	orderHandler := order.NewHandler(orderUsecase)
	customerHandler := customer.NewHandler(customerUsecase)
	paperHandler := paper.NewHandler(paperUsecase)
	phaseHandler := phase.NewHandler(phaseUsecase)
	designHandler := design.NewHandler(designUsecase)
	printHandler := print.NewHandler(printUsecase)
	finishingHandler := finishing.NewHandler(finishingUsecase)
	otherHandler := other.NewHandler(otherUsecase)
	transactionHandler := transaction.NewHandler(transactionUsecase)

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

	routerUser := router.Group("/user", checkTokenMiddleware)
	routerUser.GET("", userHandler.Page)
	routerUser.POST("", userHandler.Create)
	routerUser.POST("/change-password", userHandler.ChangePassword)
	routerUser.PUT("/:id", userHandler.Update)
	routerUser.GET("/:id", userHandler.GetById)
	routerUser.DELETE("/:id", userHandler.Delete)

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

	routerDesign := router.Group("/design", checkTokenMiddleware)
	routerDesign.GET("", designHandler.Page)
	routerDesign.POST("", designHandler.Create)
	routerDesign.PUT("/:id", designHandler.Update)
	routerDesign.GET("/:id", designHandler.GetById)
	routerDesign.DELETE("/:id", designHandler.Delete)

	routerPrint := router.Group("/print", checkTokenMiddleware)
	routerPrint.GET("", printHandler.Page)
	routerPrint.POST("", printHandler.Create)
	routerPrint.PUT("/:id", printHandler.Update)
	routerPrint.GET("/:id", printHandler.GetById)
	routerPrint.DELETE("/:id", printHandler.Delete)

	routerFinishing := router.Group("/finishing", checkTokenMiddleware)
	routerFinishing.GET("", finishingHandler.Page)
	routerFinishing.POST("", finishingHandler.Create)
	routerFinishing.PUT("/:id", finishingHandler.Update)
	routerFinishing.GET("/:id", finishingHandler.GetById)
	routerFinishing.DELETE("/:id", finishingHandler.Delete)

	routerOther := router.Group("/other", checkTokenMiddleware)
	routerOther.GET("", otherHandler.Page)
	routerOther.POST("", otherHandler.Create)
	routerOther.PUT("/:id", otherHandler.Update)
	routerOther.GET("/:id", otherHandler.GetById)
	routerOther.DELETE("/:id", otherHandler.Delete)

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

//func checkAdminMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
//	return func(c echo.Context) error {
//		userLogin := c.Get(constant.TokenUserContext).(*jwt.UserLogin)
//		if !userLogin.IsAdmin {
//			return response.ErrorForce(http.StatusForbidden, response.ErrorMiddlewareNotAdmin).SendJSON(c)
//		}
//		return next(c)
//	}
//}
