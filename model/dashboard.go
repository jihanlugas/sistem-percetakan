package model

type Dashboard struct {
	ChartTransaction    LineChart           `json:"chartTransaction"`
	TransactionOneDay   DasboardTransaction `json:"transactionOneDay"`
	TransactionOneWeek  DasboardTransaction `json:"transactionOneWeek"`
	TransactionOneMonth DasboardTransaction `json:"transactionOneMonth"`
}

type DasboardTransaction struct {
	TotalDebitCash      int64 `json:"totalDebitCash"`
	TotalKreditCash     int64 `json:"totalKreditCash"`
	TotalDebitTransfer  int64 `json:"totalDebitTransfer"`
	TotalKreditTransfer int64 `json:"totalKreditTransfer"`
	TotalOrder          int64 `json:"totalOrder"`
}

type LineChart struct {
	Label    []string  `json:"labels"`
	Datasets []Dataset `json:"datasets"`
}

type Dataset struct {
	Label           string  `json:"label"`
	Data            []int64 `json:"data"`
	BorderColor     string  `json:"borderColor"`
	BackgroundColor string  `json:"backgroundColor"`
}
