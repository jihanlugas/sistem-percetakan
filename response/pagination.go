package response

import (
	"github.com/jihanlugas/sistem-percetakan/request"
	"math"
)

type Pagination struct {
	Page        int         `json:"page"`
	DataPerPage int         `json:"dataPerPage"`
	TotalData   int64       `json:"totalData"`
	TotalPage   int         `json:"totalPage"`
	List        interface{} `json:"list" swaggertype:"array,object"`
}

func PayloadPagination(req request.IPaging, list interface{}, totalData int64) *Pagination {
	pgn := Pagination{
		Page:        req.GetPage(),
		DataPerPage: req.GetLimit(),
		TotalData:   totalData,
		TotalPage:   int(math.Ceil(float64(totalData) / float64(req.GetLimit()))),
		List:        list,
	}

	req.SetPage(0)

	return &pgn
}
