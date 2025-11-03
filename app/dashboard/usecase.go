package dashboard

import (
	"time"

	"github.com/jihanlugas/sistem-percetakan/app/order"
	"github.com/jihanlugas/sistem-percetakan/app/transaction"
	"github.com/jihanlugas/sistem-percetakan/constant"
	"github.com/jihanlugas/sistem-percetakan/db"
	"github.com/jihanlugas/sistem-percetakan/jwt"
	"github.com/jihanlugas/sistem-percetakan/model"
	"github.com/jihanlugas/sistem-percetakan/request"
	"gorm.io/gorm"
)

type Usecase interface {
	GetDashboard(loginUser jwt.UserLogin, companyID string) (dashboard model.Dashboard, err error)
}

type usecase struct {
	repositoryOrder       order.Repository
	repositoryTransaction transaction.Repository
}

func (u usecase) GetDashboard(loginUser jwt.UserLogin, companyID string) (dashboard model.Dashboard, err error) {
	var lineChart model.LineChart
	var dataDebit []int64
	var dataKredit []int64

	conn, closeConn := db.GetConnection()
	defer closeConn()

	now := time.Now()
	startDtADay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	startDtAWeek := startDtADay.AddDate(0, 0, -7)
	startDtAMonth := startDtADay.AddDate(0, 0, -30)
	endDt := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, time.Local)

	debitTransactions, err := u.repositoryTransaction.GetDailyAmountPeriod(conn, companyID, constant.TRANSACTION_TYPE_DEBIT, startDtAMonth, endDt)
	if err != nil {
		return
	}

	kreditTransactions, err := u.repositoryTransaction.GetDailyAmountPeriod(conn, companyID, constant.TRANSACTION_TYPE_KREDIT, startDtAMonth, endDt)
	if err != nil {
		return
	}

	for _, debitTransaction := range debitTransactions {
		lineChart.Label = append(lineChart.Label, debitTransaction.Date.Format(constant.FormatDateLayout))
	}

	for _, debitTransaction := range debitTransactions {
		dataDebit = append(dataDebit, debitTransaction.Amount)
	}

	for _, kreditTransaction := range kreditTransactions {
		dataKredit = append(dataKredit, kreditTransaction.Amount)
	}

	lineChart.Datasets = append(lineChart.Datasets, model.Dataset{
		Label:           "Pemasukan",
		Data:            dataDebit,
		BorderColor:     "rgb(0, 201, 81)",
		BackgroundColor: "rgba(0, 201, 81, 0.2)",
	})

	lineChart.Datasets = append(lineChart.Datasets, model.Dataset{
		Label:           "Pengeluaran",
		Data:            dataKredit,
		BorderColor:     "rgb(255, 32, 86)",
		BackgroundColor: "rgba(255, 32, 86, 0.2)",
	})

	totalDebitCashOneDay, totalKreditCashOneDay, totalDebitTransferOneDay, totalKreditTransferOneDay, totalOrderOneDay, err := u.getDataPeroid(conn, companyID, startDtADay, endDt)
	if err != nil {
		return dashboard, err
	}

	totalDebitCashOneWeek, totalKreditCashOneWeek, totalDebitTransferOneWeek, totalKreditTransferOneWeek, totalOrderOneWeek, err := u.getDataPeroid(conn, companyID, startDtAWeek, endDt)
	if err != nil {
		return dashboard, err
	}

	totalDebitCashOneMonth, totalKreditCashOneMonth, totalDebitTransferOneMonth, totalKreditTransferOneMonth, totalOrderOneMonth, err := u.getDataPeroid(conn, companyID, startDtAMonth, endDt)
	if err != nil {
		return dashboard, err
	}

	dashboard.ChartTransaction = lineChart
	dashboard.TransactionOneDay = model.DasboardTransaction{
		TotalDebitCash:      totalDebitCashOneDay,
		TotalKreditCash:     totalKreditCashOneDay,
		TotalDebitTransfer:  totalDebitTransferOneDay,
		TotalKreditTransfer: totalKreditTransferOneDay,
		TotalOrder:          totalOrderOneDay,
	}
	dashboard.TransactionOneWeek = model.DasboardTransaction{
		TotalDebitCash:      totalDebitCashOneWeek,
		TotalKreditCash:     totalKreditCashOneWeek,
		TotalDebitTransfer:  totalDebitTransferOneWeek,
		TotalKreditTransfer: totalKreditTransferOneWeek,
		TotalOrder:          totalOrderOneWeek,
	}
	dashboard.TransactionOneMonth = model.DasboardTransaction{
		TotalDebitCash:      totalDebitCashOneMonth,
		TotalKreditCash:     totalKreditCashOneMonth,
		TotalDebitTransfer:  totalDebitTransferOneMonth,
		TotalKreditTransfer: totalKreditTransferOneMonth,
		TotalOrder:          totalOrderOneMonth,
	}

	return dashboard, err
}

func (u usecase) getDataPeroid(conn *gorm.DB, companyID string, startDt, endDt time.Time) (totalDebitCash, totalKreditCash, totalDebitTrasnfer, totalKreditTrasnfer, totalOrder int64, err error) {
	totalDebitCash, err = u.repositoryTransaction.GetTotalAmountPeriodPaymentType(conn, companyID, constant.TRANSACTION_TYPE_DEBIT, constant.PAYMENT_TYPE_CASH, startDt, endDt)
	if err != nil {
		return
	}

	totalKreditCash, err = u.repositoryTransaction.GetTotalAmountPeriodPaymentType(conn, companyID, constant.TRANSACTION_TYPE_KREDIT, constant.PAYMENT_TYPE_CASH, startDt, endDt)
	if err != nil {
		return
	}

	totalDebitTrasnfer, err = u.repositoryTransaction.GetTotalAmountPeriodPaymentType(conn, companyID, constant.TRANSACTION_TYPE_DEBIT, constant.PAYMENT_TYPE_TRANSFER, startDt, endDt)
	if err != nil {
		return
	}

	totalKreditTrasnfer, err = u.repositoryTransaction.GetTotalAmountPeriodPaymentType(conn, companyID, constant.TRANSACTION_TYPE_KREDIT, constant.PAYMENT_TYPE_TRANSFER, startDt, endDt)
	if err != nil {
		return
	}

	reqOrder := request.PageOrder{
		CompanyID: companyID,
		StartDt:   &startDt,
		EndDt:     &endDt,
	}

	totalOrder, err = u.repositoryOrder.Count(conn, reqOrder)
	if err != nil {
		return
	}

	return

}

func NewUsecase(repositoryOrder order.Repository, repositoryTransaction transaction.Repository) Usecase {
	return &usecase{
		repositoryOrder:       repositoryOrder,
		repositoryTransaction: repositoryTransaction,
	}
}
