package dto

import (
	"github.com/peter-pashchenko/grpcExchangeService/internal/models"
)

type (
	ExchangeRatesDTO struct {
		Timestamp int64 `json:"timestamp"`
		Asks      []Ask `json:"asks"`
		Bids      []Bid `json:"bids"`
	}
	Ask struct {
		Price string `json:"price"`
	}
	Bid struct {
		Price string `json:"price"`
	}
)

func ConvertToExchangeRate(dto *ExchangeRatesDTO) *models.ExchangeRate {
	var res = &models.ExchangeRate{}
	res.Timestamp = dto.Timestamp
	res.AskPrice = dto.Asks[0].Price
	res.BidPrice = dto.Bids[0].Price

	return res
}
