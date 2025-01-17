package apiGarantex

import (
	"context"
	"encoding/json"
	"github.com/peter-pashchenko/grpcExchangeService/internal/application/dto"
	"go.uber.org/zap"
	"net/http"
)

const (
	urlAPI = "https://garantex.org/api/v2/depth?market="
)

type Struct struct {
	idMarket string

	logger *zap.Logger
}

func New(idMarket string, logger *zap.Logger) *Struct {
	return &Struct{idMarket, logger}
}

func (s *Struct) Call(ctx context.Context) (*dto.ExchangeRatesDTO, error) {
	url := urlAPI + s.idMarket

	response, err := http.Get(url)

	if err != nil {
		s.logger.Error(
			"method get error",
			zap.String("url", url),
			zap.Error(err),
		)
		return nil, err
	}
	defer response.Body.Close()

	var result dto.ExchangeRatesDTO

	err = json.NewDecoder(response.Body).Decode(&result)

	if err != nil {
		s.logger.Error(
			"body decoder error",
			zap.Error(err),
		)
		return nil, err
	}

	s.logger.Debug("answer from API decoded")

	return &result, nil
}
