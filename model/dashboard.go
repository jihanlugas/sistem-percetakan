package model

type Dashboard struct {
	ChartTransaction LineChart `json:"chartTransaction"`
	TotalDebit       int64     `json:"totalDebit"`
	TotalKredit      int64     `json:"totalKredit"`
	TotalOrder       int64     `json:"totalOrder"`
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
