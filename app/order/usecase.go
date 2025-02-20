package order

import (
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
	"github.com/jung-kurt/gofpdf"
	"gorm.io/gorm"
)

type Usecase interface {
	Page(loginUser jwt.UserLogin, req request.PageOrder) (vOrders []model.OrderView, count int64, err error)
	GetById(loginUser jwt.UserLogin, id string, preloads ...string) (vOrder model.OrderView, err error)
	Create(loginUser jwt.UserLogin, req request.CreateOrder) error
	Update(loginUser jwt.UserLogin, id string, req request.UpdateOrder) error
	AddPhase(loginUser jwt.UserLogin, id string, req request.AddPhase) error
	AddTransaction(loginUser jwt.UserLogin, id string, req request.AddTransaction) error
	Delete(loginUser jwt.UserLogin, id string) error
	GenerateSpk(id string) (pdf *gofpdf.Fpdf, vOrder model.OrderView, err error)
	GenerateInvoice(id string) (pdf *gofpdf.Fpdf, vOrder model.OrderView, err error)
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
		IsDone:      false,
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
	tOrder.IsDone = req.IsDone
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

func (u usecase) GenerateSpk(id string) (pdf *gofpdf.Fpdf, vOrder model.OrderView, err error) {
	conn, closeConn := db.GetConnection()
	defer closeConn()

	preloads := []string{"Designs", "Finishings", "Prints", "Prints.Paper", "Others", "Customer"}
	vOrder, err = u.repository.GetViewById(conn, id, preloads...)
	if err != nil {
		return pdf, vOrder, errors.New(fmt.Sprint("failed to get order: ", err))
	}

	pdf = u.generateSpk(vOrder)

	return pdf, vOrder, err
}

func (u usecase) generateSpkDesign(pdf *gofpdf.Fpdf, vDesigns []model.DesignView) {
	// Header dan data dari parameter
	const (
		marginH = 15.0
		lineHt  = 5.5
		cellGap = 2.0
	)
	// var colStrList [colCount]string
	type cellType struct {
		str  string
		list [][]byte
		ht   float64
	}
	type headerType struct {
		value string
		width float64
		align string
	}
	type dataType struct {
		value string
	}

	headerData := []headerType{
		{"No", 10, "C"},
		{"Nama", 60, ""},
		{"Description", 70, ""},
		{"Qty", 40, "R"},
	} // total width 180

	listData := [][]dataType{}
	for i, data := range vDesigns {
		newData := []dataType{
			{fmt.Sprintf("%d", i+1)},
			{data.Name},
			{data.Description},
			{fmt.Sprintf("%d", i+1)},
		}
		listData = append(listData, newData)
	}

	pdf.SetFont("Arial", "", 12)
	pdf.Cell(180, 10, "Design")
	pdf.Ln(10)

	pdf.SetFont("Arial", "B", 10)
	// Headers
	pdf.SetTextColor(224, 224, 224)
	pdf.SetFillColor(128, 128, 128)
	for _, data := range headerData {
		pdf.CellFormat(data.width, 10, data.value, "1", 0, "CM", true, 0, "")
	}
	pdf.Ln(-1)
	pdf.SetTextColor(24, 24, 24)
	pdf.SetFillColor(255, 255, 255)

	// Rows
	y := pdf.GetY()
	pdf.SetFont("Arial", "", 10)
	count := 0
	for _, dataRow := range listData {
		var cellList []cellType
		maxHt := lineHt
		// Cell height calculation loop
		for i, dataCell := range dataRow {
			var cell cellType
			count++
			if count > len(dataCell.value) {
				count = 1
			}
			cell.str = dataCell.value
			cell.list = pdf.SplitLines([]byte(cell.str), headerData[i].width-cellGap-cellGap)
			cell.ht = float64(len(cell.list)) * lineHt
			if cell.ht > maxHt {
				maxHt = cell.ht
			}
			cellList = append(cellList, cell)
		}

		// Cell render loop
		x := marginH
		for i := range dataRow {
			pdf.Rect(x, y, headerData[i].width, maxHt+cellGap+cellGap, "D")
			cell := cellList[i]
			cellY := y + cellGap + (maxHt-cell.ht)/2
			for splitJ := 0; splitJ < len(cell.list); splitJ++ {
				pdf.SetXY(x+cellGap, cellY)
				pdf.CellFormat(headerData[i].width-cellGap-cellGap, lineHt, string(cell.list[splitJ]), "", 0, headerData[i].align, false, 0, "")
				cellY += lineHt
			}
			x += headerData[i].width
		}
		y += maxHt + cellGap + cellGap
	}
	pdf.SetXY(marginH, y+5)
}

func (u usecase) generateSpkPrint(pdf *gofpdf.Fpdf, vPrints []model.PrintView) {
	// Header dan data dari parameter
	const (
		marginH = 15.0
		lineHt  = 5.5
		cellGap = 2.0
	)
	// var colStrList [colCount]string
	type cellType struct {
		str  string
		list [][]byte
		ht   float64
	}
	type headerType struct {
		value string
		width float64
		align string
	}
	type dataType struct {
		value string
	}

	headerData := []headerType{
		{"No", 10, "C"},
		{"Nama", 40, ""},
		{"Description", 45, ""},
		{"Kertas", 40, ""},
		{"2 Muka", 15, ""},
		{"Lembar", 15, "R"},
		{"Qty", 15, "R"},
	} // total width 180

	listData := [][]dataType{}
	for i, data := range vPrints {
		newData := []dataType{
			{fmt.Sprintf("%d", i+1)},
			{data.Name},
			{data.Description},
			{data.Paper.Name},
			{utils.DisplayBool(data.IsDuplex, "Ya", "Tidak")},
			{fmt.Sprintf("%d", data.PageCount)},
			{fmt.Sprintf("%d", data.Qty)},
		}
		listData = append(listData, newData)
	}

	pdf.SetFont("Arial", "", 12)
	pdf.Cell(180, 10, "Print")
	pdf.Ln(10)

	pdf.SetFont("Arial", "B", 10)
	// Headers
	pdf.SetTextColor(224, 224, 224)
	pdf.SetFillColor(128, 128, 128)
	for _, data := range headerData {
		pdf.CellFormat(data.width, 10, data.value, "1", 0, "CM", true, 0, "")
	}
	pdf.Ln(-1)
	pdf.SetTextColor(24, 24, 24)
	pdf.SetFillColor(255, 255, 255)

	// Rows
	y := pdf.GetY()
	pdf.SetFont("Arial", "", 10)
	count := 0
	for _, dataRow := range listData {
		var cellList []cellType
		maxHt := lineHt
		// Cell height calculation loop
		for i, dataCell := range dataRow {
			var cell cellType
			count++
			if count > len(dataCell.value) {
				count = 1
			}
			cell.str = dataCell.value
			cell.list = pdf.SplitLines([]byte(cell.str), headerData[i].width-cellGap-cellGap)
			cell.ht = float64(len(cell.list)) * lineHt
			if cell.ht > maxHt {
				maxHt = cell.ht
			}
			cellList = append(cellList, cell)
		}

		// Cell render loop
		x := marginH
		for i := range dataRow {
			pdf.Rect(x, y, headerData[i].width, maxHt+cellGap+cellGap, "D")
			cell := cellList[i]
			cellY := y + cellGap + (maxHt-cell.ht)/2
			for splitJ := 0; splitJ < len(cell.list); splitJ++ {
				pdf.SetXY(x+cellGap, cellY)
				pdf.CellFormat(headerData[i].width-cellGap-cellGap, lineHt, string(cell.list[splitJ]), "", 0, headerData[i].align, false, 0, "")
				cellY += lineHt
			}
			x += headerData[i].width
		}
		y += maxHt + cellGap + cellGap
	}
	pdf.SetXY(marginH, y+5)
}

func (u usecase) generateSpkFinishing(pdf *gofpdf.Fpdf, vFinishings []model.FinishingView) {
	// Header dan data dari parameter
	const (
		marginH = 15.0
		lineHt  = 5.5
		cellGap = 2.0
	)
	// var colStrList [colCount]string
	type cellType struct {
		str  string
		list [][]byte
		ht   float64
	}
	type headerType struct {
		value string
		width float64
		align string
	}
	type dataType struct {
		value string
	}

	headerData := []headerType{
		{"No", 10, "C"},
		{"Nama", 60, ""},
		{"Description", 70, ""},
		{"Qty", 40, "R"},
	} // total width 180

	listData := [][]dataType{}
	for i, data := range vFinishings {
		newData := []dataType{
			{fmt.Sprintf("%d", i+1)},
			{data.Name},
			{data.Description},
			{fmt.Sprintf("%d", i+1)},
		}
		listData = append(listData, newData)
	}

	pdf.SetFont("Arial", "", 12)
	pdf.Cell(180, 10, "Finishing")
	pdf.Ln(10)

	pdf.SetFont("Arial", "B", 10)
	// Headers
	pdf.SetTextColor(224, 224, 224)
	pdf.SetFillColor(128, 128, 128)
	for _, data := range headerData {
		pdf.CellFormat(data.width, 10, data.value, "1", 0, "CM", true, 0, "")
	}
	pdf.Ln(-1)
	pdf.SetTextColor(24, 24, 24)
	pdf.SetFillColor(255, 255, 255)

	// Rows
	y := pdf.GetY()
	pdf.SetFont("Arial", "", 10)
	count := 0
	for _, dataRow := range listData {
		var cellList []cellType
		maxHt := lineHt
		// Cell height calculation loop
		for i, dataCell := range dataRow {
			var cell cellType
			count++
			if count > len(dataCell.value) {
				count = 1
			}
			cell.str = dataCell.value
			cell.list = pdf.SplitLines([]byte(cell.str), headerData[i].width-cellGap-cellGap)
			cell.ht = float64(len(cell.list)) * lineHt
			if cell.ht > maxHt {
				maxHt = cell.ht
			}
			cellList = append(cellList, cell)
		}

		// Cell render loop
		x := marginH
		for i := range dataRow {
			pdf.Rect(x, y, headerData[i].width, maxHt+cellGap+cellGap, "D")
			cell := cellList[i]
			cellY := y + cellGap + (maxHt-cell.ht)/2
			for splitJ := 0; splitJ < len(cell.list); splitJ++ {
				pdf.SetXY(x+cellGap, cellY)
				pdf.CellFormat(headerData[i].width-cellGap-cellGap, lineHt, string(cell.list[splitJ]), "", 0, headerData[i].align, false, 0, "")
				cellY += lineHt
			}
			x += headerData[i].width
		}
		y += maxHt + cellGap + cellGap
	}
	pdf.SetXY(marginH, y+5)
}

func (u usecase) generateSpkOther(pdf *gofpdf.Fpdf, vOthers []model.OtherView) {
	// Header dan data dari parameter
	const (
		marginH = 15.0
		lineHt  = 5.5
		cellGap = 2.0
	)
	// var colStrList [colCount]string
	type cellType struct {
		str  string
		list [][]byte
		ht   float64
	}
	type headerType struct {
		value string
		width float64
		align string
	}
	type dataType struct {
		value string
	}

	headerData := []headerType{
		{"No", 10, "C"},
		{"Nama", 60, ""},
		{"Description", 70, ""},
		{"Qty", 40, "R"},
	} // total width 180

	listData := [][]dataType{}
	for i, data := range vOthers {
		newData := []dataType{
			{fmt.Sprintf("%d", i+1)},
			{data.Name},
			{data.Description},
			{fmt.Sprintf("%d", i+1)},
		}
		listData = append(listData, newData)
	}

	pdf.SetFont("Arial", "", 12)
	pdf.Cell(180, 10, "Other")
	pdf.Ln(10)

	pdf.SetFont("Arial", "B", 10)
	// Headers
	pdf.SetTextColor(224, 224, 224)
	pdf.SetFillColor(128, 128, 128)
	for _, data := range headerData {
		pdf.CellFormat(data.width, 10, data.value, "1", 0, "CM", true, 0, "")
	}
	pdf.Ln(-1)
	pdf.SetTextColor(24, 24, 24)
	pdf.SetFillColor(255, 255, 255)

	// Rows
	y := pdf.GetY()
	pdf.SetFont("Arial", "", 10)
	count := 0
	for _, dataRow := range listData {
		var cellList []cellType
		maxHt := lineHt
		// Cell height calculation loop
		for i, dataCell := range dataRow {
			var cell cellType
			count++
			if count > len(dataCell.value) {
				count = 1
			}
			cell.str = dataCell.value
			cell.list = pdf.SplitLines([]byte(cell.str), headerData[i].width-cellGap-cellGap)
			cell.ht = float64(len(cell.list)) * lineHt
			if cell.ht > maxHt {
				maxHt = cell.ht
			}
			cellList = append(cellList, cell)
		}

		// Cell render loop
		x := marginH
		for i := range dataRow {
			pdf.Rect(x, y, headerData[i].width, maxHt+cellGap+cellGap, "D")
			cell := cellList[i]
			cellY := y + cellGap + (maxHt-cell.ht)/2
			for splitJ := 0; splitJ < len(cell.list); splitJ++ {
				pdf.SetXY(x+cellGap, cellY)
				pdf.CellFormat(headerData[i].width-cellGap-cellGap, lineHt, string(cell.list[splitJ]), "", 0, headerData[i].align, false, 0, "")
				cellY += lineHt
			}
			x += headerData[i].width
		}
		y += maxHt + cellGap + cellGap
	}
	pdf.SetXY(marginH, y+5)
}

func (u usecase) generateSpk(vOrder model.OrderView) (pdf *gofpdf.Fpdf) {

	pdf = gofpdf.New("P", "mm", "A4", "") // 210 x 297
	pdf.SetMargins(15, 15, 15)
	pdf.AddPage()

	// Judul Surat Perintah Kerja
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(190, 10, "Surat Perintah Kerja")
	pdf.Ln(12)

	// Informasi pekerjaan
	pdf.SetFont("Arial", "", 12)
	pdf.MultiCell(0, 10, "Berikut adalah rincian perintah kerja untuk pencetakan buku dalam bentuk tabel:", "", "", false)

	if len(vOrder.Designs) > 0 {
		u.generateSpkDesign(pdf, vOrder.Designs)
	}
	if len(vOrder.Prints) > 0 {
		u.generateSpkPrint(pdf, vOrder.Prints)
	}
	if len(vOrder.Finishings) > 0 {
		u.generateSpkFinishing(pdf, vOrder.Finishings)
	}
	if len(vOrder.Others) > 0 {
		u.generateSpkOther(pdf, vOrder.Others)
	}

	// Tanda Tangan
	//pdf.Ln(20)
	//pdf.Cell(190, 10, "Hormat Kami,")
	//pdf.Ln(15)
	//pdf.Cell(190, 10, "(______________________)")

	return pdf
}

func (u usecase) GenerateInvoice(id string) (pdf *gofpdf.Fpdf, vOrder model.OrderView, err error) {
	conn, closeConn := db.GetConnection()
	defer closeConn()

	preloads := []string{"Designs", "Finishings", "Prints", "Prints.Paper", "Others", "Customer", "Transactions"}
	vOrder, err = u.repository.GetViewById(conn, id, preloads...)
	if err != nil {
		return pdf, vOrder, errors.New(fmt.Sprint("failed to get order: ", err))
	}

	pdf = u.generateInvoice(vOrder)

	return pdf, vOrder, err
}

func (u usecase) generateInvoice(vOrder model.OrderView) (pdf *gofpdf.Fpdf) {

	pdf = gofpdf.New("P", "mm", "A4", "") // 210 x 297
	pdf.SetMargins(15, 10, 15)

	totalPagesAlias := "{nb}"

	pdf.SetHeaderFunc(func() {
		// Add a header
		pdf.SetFont("Arial", "B", 12)
		pdf.Cell(0, 10, fmt.Sprintf("Custom Header - Page: %d", pdf.PageNo()))
		pdf.Ln(10) // Line break
	})

	pdf.SetFooterFunc(func() {
		// Add a footer
		pdf.SetY(-15) // Position at 1.5 cm from the bottom
		pdf.SetFont("Arial", "I", 8)
		pdf.Cell(0, 10, fmt.Sprintf("Footer - Page: %d of %s", pdf.PageNo(), totalPagesAlias))
	})

	pdf.AliasNbPages(totalPagesAlias) // Register the alias
	pdf.AddPage()
	pdf.SetFont("Arial", "", 12)
	for i := 0; i < 200; i++ {
		pdf.Cell(0, 10, "This is line "+string(i+1))
		pdf.Ln(10) // Line break
	}

	//// Header
	//pdf.SetFont("Arial", "B", 16)
	//pdf.Cell(0, 10, "Invoice")
	//pdf.Ln(12)
	//
	//// Informasi Perusahaan
	//pdf.SetFont("Arial", "", 12)
	//pdf.Cell(0, 10, fmt.Sprintf("Dari : %s", vOrder.CompanyName))
	////pdf.Ln(6)
	////pdf.Cell(0, 10, "Alamat: Jalan Merdeka No.123, Jakarta")
	////pdf.Ln(6)
	////pdf.Cell(0, 10, "Email: contoh@perusahaan.com")
	//pdf.Ln(10)
	//
	//// Informasi Pelanggan
	//pdf.SetFont("Arial", "B", 12)
	//pdf.Cell(0, 10, "Kepada:")
	//pdf.Ln(6)
	//pdf.SetFont("Arial", "", 12)
	//pdf.Cell(0, 10, fmt.Sprintf("Name : %s", vOrder.Customer.Name))
	////pdf.Ln(6)
	////pdf.Cell(0, 10, "Alamat: Jalan Sudirman No.456, Bandung")
	//pdf.Ln(10)
	//
	//// Daftar Item
	//pdf.SetFont("Arial", "B", 12)
	//pdf.Cell(90, 10, "Deskripsi")
	//pdf.Cell(30, 10, "Jumlah")
	//pdf.Cell(30, 10, "Harga Satuan")
	//pdf.Cell(30, 10, "Total")
	//pdf.Ln(10)
	//
	//items := []struct {
	//	Description string
	//	Quantity    int
	//	UnitPrice   float64
	//}{
	//	{"Item 1", 2, 50000},
	//	{"Item 2", 1, 75000},
	//	{"Item 3", 3, 30000},
	//}
	//
	//var grandTotal float64
	//
	//pdf.SetFont("Arial", "", 12)
	//for _, item := range items {
	//	total := float64(item.Quantity) * item.UnitPrice
	//	grandTotal += total
	//
	//	pdf.Cell(90, 10, item.Description)
	//	pdf.CellFormat(30, 10, fmt.Sprintf("%d", item.Quantity), "", 0, "C", false, 0, "")
	//	pdf.CellFormat(30, 10, fmt.Sprintf("Rp%.2f", item.UnitPrice), "", 0, "R", false, 0, "")
	//	pdf.CellFormat(30, 10, fmt.Sprintf("Rp%.2f", total), "", 0, "R", false, 0, "")
	//	pdf.Ln(10)
	//}
	//
	//// Total Keseluruhan
	//pdf.SetFont("Arial", "B", 12)
	//pdf.Cell(150, 10, "Grand Total")
	//pdf.CellFormat(30, 10, fmt.Sprintf("Rp%.2f", grandTotal), "", 0, "R", false, 0, "")

	// Tanda Tangan
	//pdf.Ln(20)
	//pdf.Cell(190, 10, "Hormat Kami,")
	//pdf.Ln(15)
	//pdf.Cell(190, 10, "(______________________)")

	return pdf
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
