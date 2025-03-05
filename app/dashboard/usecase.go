package dashboard

import (
	"github.com/jihanlugas/sistem-percetakan/app/order"
	"github.com/jihanlugas/sistem-percetakan/app/transaction"
	"github.com/jihanlugas/sistem-percetakan/constant"
	"github.com/jihanlugas/sistem-percetakan/db"
	"github.com/jihanlugas/sistem-percetakan/jwt"
	"github.com/jihanlugas/sistem-percetakan/model"
	"github.com/jihanlugas/sistem-percetakan/request"
	"time"
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
	var totalDebit int64
	var totalKredit int64
	var totalOrder int64

	conn, closeConn := db.GetConnection()
	defer closeConn()

	now := time.Now()
	startDtADay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	startDtAWeek := startDtADay.AddDate(0, 0, -6)
	endDt := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, time.Local)

	totalDebit, err = u.repositoryTransaction.GetTotalAmountPeriod(conn, companyID, constant.TRANSACTION_TYPE_DEBIT, startDtADay, endDt)
	if err != nil {
		return
	}

	totalKredit, err = u.repositoryTransaction.GetTotalAmountPeriod(conn, companyID, constant.TRANSACTION_TYPE_KREDIT, startDtADay, endDt)
	if err != nil {
		return
	}

	debitTransactions, err := u.repositoryTransaction.GetDailyAmountPeriod(conn, companyID, constant.TRANSACTION_TYPE_DEBIT, startDtAWeek, endDt)
	if err != nil {
		return
	}

	kreditTransactions, err := u.repositoryTransaction.GetDailyAmountPeriod(conn, companyID, constant.TRANSACTION_TYPE_KREDIT, startDtAWeek, endDt)
	if err != nil {
		return
	}

	reqOrder := request.PageOrder{
		CompanyID: companyID,
		StartDt:   &startDtADay,
		EndDt:     &endDt,
	}
	totalOrder, err = u.repositoryOrder.Count(conn, reqOrder)
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

	dashboard.ChartTransaction = lineChart
	dashboard.TotalDebit = totalDebit
	dashboard.TotalKredit = totalKredit
	dashboard.TotalOrder = totalOrder

	return dashboard, err
}

func NewUsecase(repositoryOrder order.Repository, repositoryTransaction transaction.Repository) Usecase {
	return &usecase{
		repositoryOrder:       repositoryOrder,
		repositoryTransaction: repositoryTransaction,
	}
}
