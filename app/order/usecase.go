package order

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/jihanlugas/sistem-percetakan/app/customer"
	"github.com/jihanlugas/sistem-percetakan/app/design"
	"github.com/jihanlugas/sistem-percetakan/app/finishing"
	"github.com/jihanlugas/sistem-percetakan/app/orderphase"
	"github.com/jihanlugas/sistem-percetakan/app/other"
	"github.com/jihanlugas/sistem-percetakan/app/phase"
	"github.com/jihanlugas/sistem-percetakan/app/print"
	"github.com/jihanlugas/sistem-percetakan/app/transaction"
	"github.com/jihanlugas/sistem-percetakan/constant"
	"github.com/jihanlugas/sistem-percetakan/db"
	"github.com/jihanlugas/sistem-percetakan/jwt"
	"github.com/jihanlugas/sistem-percetakan/model"
	"github.com/jihanlugas/sistem-percetakan/request"
	"github.com/jihanlugas/sistem-percetakan/response"
	"github.com/jihanlugas/sistem-percetakan/utils"
	"gorm.io/gorm"
	"html/template"
	"os"
)

type Usecase interface {
	Page(loginUser jwt.UserLogin, req request.PageOrder) (vOrders []model.OrderView, count int64, err error)
	GetById(loginUser jwt.UserLogin, id string, preloads ...string) (vOrder model.OrderView, err error)
	Create(loginUser jwt.UserLogin, req request.CreateOrder) error
	Update(loginUser jwt.UserLogin, id string, req request.UpdateOrder) error
	AddPhase(loginUser jwt.UserLogin, id string, req request.AddPhase) error
	AddTransaction(loginUser jwt.UserLogin, id string, req request.AddTransaction) error
	Delete(loginUser jwt.UserLogin, id string) error
	GenerateSpk(id string) (pdfBytes []byte, vOrder model.OrderView, err error)
	GenerateInvoice(id string) (pdfBytes []byte, vOrder model.OrderView, err error)
}

type usecase struct {
	repository            Repository
	repositoryDesign      design.Repository
	repositoryPrint       print.Repository
	repositoryFinishing   finishing.Repository
	repositoryOther       other.Repository
	repositoryOrderphase  orderphase.Repository
	repositoryCustomer    customer.Repository
	repositoryPhase       phase.Repository
	repositoryTransaction transaction.Repository
}

func (u usecase) Page(loginUser jwt.UserLogin, req request.PageOrder) (vOrders []model.OrderView, count int64, err error) {
	conn, closeConn := db.GetConnection()
	defer closeConn()

	vOrders, count, err = u.repository.Page(conn, req)
	if err != nil {
		return vOrders, count, err
	}

	return vOrders, count, err
}

func (u usecase) GetById(loginUser jwt.UserLogin, id string, preloads ...string) (vOrder model.OrderView, err error) {
	conn, closeConn := db.GetConnection()
	defer closeConn()

	vOrder, err = u.repository.GetViewById(conn, id, preloads...)
	if err != nil {
		return vOrder, errors.New(fmt.Sprint("failed to get order: ", err))
	}

	if jwt.IsSaveCompanyIDOR(loginUser, vOrder.CompanyID) {
		return vOrder, errors.New(response.ErrorHandlerIDOR)
	}

	return vOrder, err
}

func (u usecase) Create(loginUser jwt.UserLogin, req request.CreateOrder) error {
	var err error
	var tOrder model.Order

	conn, closeConn := db.GetConnection()
	defer closeConn()

	tx := conn.Begin()

	if req.CustomerID == "" && req.NewCustomer != "" {
		newCustomer := model.Customer{
			ID:          utils.GetUniqueID(),
			CompanyID:   req.CompanyID,
			Name:        req.NewCustomer,
			PhoneNumber: utils.FormatPhoneTo62(req.NewCustomerPhone),
			CreateBy:    loginUser.UserID,
			UpdateBy:    loginUser.UserID,
		}

		err = u.repositoryCustomer.Create(tx, newCustomer)
		if err != nil {
			return errors.New(fmt.Sprint("failed to create customer: ", err))
		}

		req.CustomerID = newCustomer.ID
	}

	tOrder = model.Order{
		ID:          utils.GetUniqueID(),
		CompanyID:   req.CompanyID,
		CustomerID:  req.CustomerID,
		Name:        req.Name,
		Description: req.Description,
		Number:      u.repository.GetNextNumber(tx, req.CompanyID),
		CreateBy:    loginUser.UserID,
		UpdateBy:    loginUser.UserID,
	}

	err = u.repository.Create(tx, tOrder)
	if err != nil {
		return errors.New(fmt.Sprint("failed to create order: ", err))
	}

	err = u.createOrderDesigns(tx, loginUser, req.Designs, req.CompanyID, tOrder.ID)
	if err != nil {
		return errors.New(fmt.Sprint("failed to create designs: ", err))
	}

	err = u.createOrderPrints(tx, loginUser, req.Prints, req.CompanyID, tOrder.ID)
	if err != nil {
		return errors.New(fmt.Sprint("failed to create prints: ", err))
	}

	err = u.createOrderFinishings(tx, loginUser, req.Finishings, req.CompanyID, tOrder.ID)
	if err != nil {
		return errors.New(fmt.Sprint("failed to create finishings: ", err))
	}

	err = u.createOrderOthers(tx, loginUser, req.Others, req.CompanyID, tOrder.ID)
	if err != nil {
		return errors.New(fmt.Sprint("failed to create others: ", err))
	}

	err = u.createOrderOrderphase(tx, loginUser, req.OrderphaseID, req.CompanyID, tOrder.ID)
	if err != nil {
		return errors.New(fmt.Sprint("failed to create orderphase: ", err))
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return err
}

func (u usecase) Update(loginUser jwt.UserLogin, id string, req request.UpdateOrder) error {
	var err error
	var tOrder model.Order

	conn, closeConn := db.GetConnection()
	defer closeConn()

	tOrder, err = u.repository.GetTableById(conn, id)
	if err != nil {
		return errors.New(fmt.Sprint("failed to get order: ", err))
	}

	if jwt.IsSaveCompanyIDOR(loginUser, tOrder.CompanyID) {
		return errors.New(response.ErrorHandlerIDOR)
	}

	tx := conn.Begin()

	tOrder.Name = req.Name
	tOrder.Description = req.Description
	tOrder.UpdateBy = loginUser.UserID
	err = u.repository.Save(tx, tOrder)
	if err != nil {
		return errors.New(fmt.Sprint("failed to update order: ", err))
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return err
}

func (u usecase) AddPhase(loginUser jwt.UserLogin, id string, req request.AddPhase) error {
	var err error
	var tOrder model.Order
	var tPhase model.Phase
	var tOrderphase model.Orderphase

	conn, closeConn := db.GetConnection()
	defer closeConn()

	tOrder, err = u.repository.GetTableById(conn, id)
	if err != nil {
		return errors.New(fmt.Sprint("failed to get order: ", err))
	}

	if jwt.IsSaveCompanyIDOR(loginUser, tOrder.CompanyID) {
		return errors.New(response.ErrorHandlerIDOR)
	}

	tPhase, err = u.repositoryPhase.GetTableById(conn, req.OrderphaseID)
	if err != nil {
		return errors.New(fmt.Sprint("failed to get phase: ", err))
	}

	if jwt.IsSaveCompanyIDOR(loginUser, tPhase.CompanyID) {
		return errors.New(response.ErrorHandlerIDOR)
	}

	tx := conn.Begin()

	tOrderphase = model.Orderphase{
		ID:        utils.GetUniqueID(),
		CompanyID: loginUser.CompanyID,
		OrderID:   tOrder.ID,
		PhaseID:   tPhase.ID,
		Name:      tPhase.Name,
		CreateBy:  loginUser.UserID,
		UpdateBy:  loginUser.UserID,
	}
	err = u.repositoryOrderphase.Create(tx, tOrderphase)
	if err != nil {
		return err
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return err
}

func (u usecase) AddTransaction(loginUser jwt.UserLogin, id string, req request.AddTransaction) error {
	var err error
	var tOrder model.Order
	var tTransaction model.Transaction

	conn, closeConn := db.GetConnection()
	defer closeConn()

	tOrder, err = u.repository.GetTableById(conn, id)
	if err != nil {
		return errors.New(fmt.Sprint("failed to get order: ", err))
	}

	if jwt.IsSaveCompanyIDOR(loginUser, tOrder.CompanyID) {
		return errors.New(response.ErrorHandlerIDOR)
	}

	tx := conn.Begin()

	tTransaction = model.Transaction{
		ID:          utils.GetUniqueID(),
		CompanyID:   loginUser.CompanyID,
		OrderID:     tOrder.ID,
		Name:        req.Name,
		Amount:      req.Amount,
		Description: req.Description,
		Type:        constant.TRANSACTION_TYPE_DEBIT,
		CreateBy:    loginUser.UserID,
		UpdateBy:    loginUser.UserID,
	}
	err = u.repositoryTransaction.Create(tx, tTransaction)
	if err != nil {
		return errors.New(fmt.Sprint("failed to create transaction: ", err))
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return err
}

func (u usecase) createOrderDesigns(tx *gorm.DB, loginUser jwt.UserLogin, req []request.CreateOrderDesign, companyID, orderId string) error {
	var err error
	var tDesigns []model.Design
	if len(req) > 0 {
		for _, reqDesign := range req {
			tDesign := model.Design{
				ID:          utils.GetUniqueID(),
				CompanyID:   companyID,
				OrderID:     orderId,
				Name:        reqDesign.Name,
				Description: reqDesign.Description,
				Qty:         reqDesign.Qty,
				Price:       reqDesign.Price,
				Total:       reqDesign.Total,
				CreateBy:    loginUser.UserID,
				UpdateBy:    loginUser.UserID,
			}
			tDesigns = append(tDesigns, tDesign)
		}

		err = u.repositoryDesign.Creates(tx, tDesigns)
		if err != nil {
			return err
		}
	}

	return err
}

func (u usecase) createOrderPrints(tx *gorm.DB, loginUser jwt.UserLogin, req []request.CreateOrderPrint, companyId, orderId string) error {
	var err error
	var tPrints []model.Print
	if len(req) > 0 {
		for _, reqPrint := range req {
			tPrint := model.Print{
				ID:          utils.GetUniqueID(),
				CompanyID:   companyId,
				OrderID:     orderId,
				PaperID:     reqPrint.PaperID,
				Name:        reqPrint.Name,
				Description: reqPrint.Description,
				IsDuplex:    reqPrint.IsDuplex,
				PageCount:   reqPrint.PageCount,
				Qty:         reqPrint.Qty,
				Price:       reqPrint.Price,
				Total:       reqPrint.Total,
				CreateBy:    loginUser.UserID,
				UpdateBy:    loginUser.UserID,
			}
			tPrints = append(tPrints, tPrint)
		}

		err = u.repositoryPrint.Creates(tx, tPrints)
		if err != nil {
			return err
		}
	}
	return err
}

func (u usecase) createOrderFinishings(tx *gorm.DB, loginUser jwt.UserLogin, req []request.CreateOrderFinishing, companyId, orderId string) error {
	var err error
	var tFinishings []model.Finishing
	if len(req) > 0 {
		for _, reqFinishing := range req {
			tFinishing := model.Finishing{
				ID:          utils.GetUniqueID(),
				CompanyID:   companyId,
				OrderID:     orderId,
				Name:        reqFinishing.Name,
				Description: reqFinishing.Description,
				Qty:         reqFinishing.Qty,
				Price:       reqFinishing.Price,
				Total:       reqFinishing.Total,
				CreateBy:    loginUser.UserID,
				UpdateBy:    loginUser.UserID,
			}
			tFinishings = append(tFinishings, tFinishing)
		}

		err = u.repositoryFinishing.Creates(tx, tFinishings)
		if err != nil {
			return err
		}
	}
	return err
}

func (u usecase) createOrderOthers(tx *gorm.DB, loginUser jwt.UserLogin, req []request.CreateOrderOther, companyId, orderId string) error {
	var err error
	var tOthers []model.Other
	if len(req) > 0 {
		for _, reqOther := range req {
			tOther := model.Other{
				ID:          utils.GetUniqueID(),
				CompanyID:   companyId,
				OrderID:     orderId,
				Name:        reqOther.Name,
				Description: reqOther.Description,
				Qty:         reqOther.Qty,
				Price:       reqOther.Price,
				Total:       reqOther.Total,
				CreateBy:    loginUser.UserID,
				UpdateBy:    loginUser.UserID,
			}
			tOthers = append(tOthers, tOther)
		}

		err = u.repositoryOther.Creates(tx, tOthers)
		if err != nil {
			return err
		}
	}
	return err
}

func (u usecase) createOrderOrderphase(tx *gorm.DB, loginUser jwt.UserLogin, req, companyId, orderId string) error {
	var err error
	var tOrderphase model.Orderphase
	var tPhase model.Phase

	if req != "" {
		tPhase, err = u.repositoryPhase.GetTableById(tx, req)
		if err != nil {
			return errors.New(fmt.Sprint("failed to get phase: ", err))
		}

		tOrderphase = model.Orderphase{
			ID:        utils.GetUniqueID(),
			CompanyID: companyId,
			OrderID:   orderId,
			PhaseID:   tPhase.ID,
			Name:      tPhase.Name,
			CreateBy:  loginUser.UserID,
			UpdateBy:  loginUser.UserID,
		}
		err = u.repositoryOrderphase.Create(tx, tOrderphase)
		if err != nil {
			return err
		}
	}

	return err
}

func (u usecase) Delete(loginUser jwt.UserLogin, id string) error {
	var err error
	var tOrder model.Order

	conn, closeConn := db.GetConnection()
	defer closeConn()

	tOrder, err = u.repository.GetTableById(conn, id)
	if err != nil {
		return errors.New(fmt.Sprint("failed to get order: ", err))
	}

	if jwt.IsSaveCompanyIDOR(loginUser, tOrder.CompanyID) {
		return errors.New(response.ErrorHandlerIDOR)
	}

	tx := conn.Begin()

	err = u.repository.Delete(tx, tOrder)
	if err != nil {
		return errors.New(fmt.Sprint("failed to delete order: ", err))
	}

	err = u.repositoryDesign.DeleteByOrderId(tx, tOrder.ID)
	if err != nil {
		return errors.New(fmt.Sprint("failed to delete order design: ", err))
	}

	err = u.repositoryFinishing.DeleteByOrderId(tx, tOrder.ID)
	if err != nil {
		return errors.New(fmt.Sprint("failed to delete order finishing: ", err))
	}

	err = u.repositoryPrint.DeleteByOrderId(tx, tOrder.ID)
	if err != nil {
		return errors.New(fmt.Sprint("failed to delete order print: ", err))
	}

	err = u.repositoryOther.DeleteByOrderId(tx, tOrder.ID)
	if err != nil {
		return errors.New(fmt.Sprint("failed to delete order other: ", err))
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return err
}

func (u usecase) GenerateSpk(id string) (pdfBytes []byte, vOrder model.OrderView, err error) {
	conn, closeConn := db.GetConnection()
	defer closeConn()

	preloads := []string{"Company", "Designs", "Finishings", "Prints", "Prints.Paper", "Others", "Customer"}
	vOrder, err = u.repository.GetViewById(conn, id, preloads...)
	if err != nil {
		return pdfBytes, vOrder, errors.New(fmt.Sprint("failed to get order: ", err))
	}

	pdfBytes, err = u.generateSpk(vOrder)

	return pdfBytes, vOrder, err
}

func (u usecase) generateSpk(vOrder model.OrderView) (pdfBytes []byte, err error) {
	tmpl := template.New("spk.html").Funcs(template.FuncMap{
		"displayLembar": func(lembar int64) string {
			return fmt.Sprintf("%s Lembar", utils.DisplayNumber(lembar))
		},
		"displayDuplex": func(isDuplex bool) string {
			if isDuplex {
				return "2 Muka"
			}
			return "1 Muka"
		},
		"displayDate":        utils.DisplayDate,
		"displayDatetime":    utils.DisplayDatetime,
		"displayNumber":      utils.DisplayNumber,
		"displayMoney":       utils.DisplayMoney,
		"displayPhoneNumber": utils.DisplayPhoneNumber,
		"displaySpkNumber":   utils.DisplaySpkNumber,
	})

	// Parse template setelah fungsi didaftarkan
	tmpl, err = tmpl.ParseFiles("assets/template/spk.html")
	if err != nil {
		return pdfBytes, err
	}

	// Render template ke buffer
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, vOrder); err != nil {
		return pdfBytes, err
	}

	// Simpan HTML render ke file sementara
	tempHTMLFile := "temp.html"
	if err := os.WriteFile(tempHTMLFile, buf.Bytes(), 0644); err != nil {
		return pdfBytes, err
	}
	defer os.Remove(tempHTMLFile)

	return utils.GeneratePDFWithChromedp(tempHTMLFile)
}

func (u usecase) GenerateInvoice(id string) (pdfBytes []byte, vOrder model.OrderView, err error) {
	conn, closeConn := db.GetConnection()
	defer closeConn()

	preloads := []string{"Company", "Designs", "Finishings", "Prints", "Prints.Paper", "Others", "Customer", "Transactions"}
	vOrder, err = u.repository.GetViewById(conn, id, preloads...)
	if err != nil {
		return pdfBytes, vOrder, errors.New(fmt.Sprint("failed to get order: ", err))
	}

	pdfBytes, err = u.generateInvoice(vOrder)

	return pdfBytes, vOrder, err
}

func (u usecase) generateInvoice(vOrder model.OrderView) (pdfBytes []byte, err error) {
	tmpl := template.New("invoice.html").Funcs(template.FuncMap{
		"displayLembar": func(lembar int64) string {
			return fmt.Sprintf("%s Lembar", utils.DisplayNumber(lembar))
		},
		"displayDuplex": func(isDuplex bool) string {
			if isDuplex {
				return "2 Muka"
			}
			return "1 Muka"
		},
		"displayDate":          utils.DisplayDate,
		"displayDatetime":      utils.DisplayDatetime,
		"displayNumber":        utils.DisplayNumber,
		"displayMoney":         utils.DisplayMoney,
		"displayPhoneNumber":   utils.DisplayPhoneNumber,
		"displayInvoiceNumber": utils.DisplayInvoiceNumber,
	})

	// Parse template setelah fungsi didaftarkan
	tmpl, err = tmpl.ParseFiles("assets/template/invoice.html")
	if err != nil {
		return pdfBytes, err
	}

	// Render template ke buffer
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, vOrder); err != nil {
		return pdfBytes, err
	}

	// Simpan HTML render ke file sementara
	tempHTMLFile := "temp.html"
	if err := os.WriteFile(tempHTMLFile, buf.Bytes(), 0644); err != nil {
		return pdfBytes, err
	}
	defer os.Remove(tempHTMLFile)

	return utils.GeneratePDFWithChromedp(tempHTMLFile)
}

func NewUsecase(repository Repository, repositoryDesign design.Repository, repositoryPrint print.Repository, repositoryFinishing finishing.Repository, repositoryOther other.Repository, repositoryOrderphase orderphase.Repository, repositoryCustomer customer.Repository, repositoryPhase phase.Repository, repositoryTransaction transaction.Repository) Usecase {
	return &usecase{
		repository:            repository,
		repositoryDesign:      repositoryDesign,
		repositoryPrint:       repositoryPrint,
		repositoryFinishing:   repositoryFinishing,
		repositoryOther:       repositoryOther,
		repositoryOrderphase:  repositoryOrderphase,
		repositoryCustomer:    repositoryCustomer,
		repositoryPhase:       repositoryPhase,
		repositoryTransaction: repositoryTransaction,
	}
}
