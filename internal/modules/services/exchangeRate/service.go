package exchangeRatesService

import (
	"context"
	"github.com/peter-pashchenko/grpcExchangeService/internal/application/dto"
	"github.com/peter-pashchenko/grpcExchangeService/internal/models"
	"go.uber.org/zap"
)

type Service struct {
	repo Repo

	api Api

	logger *zap.Logger
}

//go:generate mockery --all --case underscore  --with-expecter
type Repo interface {
	SaveRate(ctx context.Context, rate *models.ExchangeRate) error
	CheckConnection(ctx context.Context) error
}

type Api interface {
	Call(ctx context.Context) (*dto.ExchangeRatesDTO, error)
}

func New(repo Repo, api Api, logger *zap.Logger) *Service {
	return &Service{repo: repo, api: api, logger: logger}
}

func (s Service) GetRate(ctx context.Context) (*models.ExchangeRate, error) {
	result, err := s.api.Call(ctx)

	if err != nil {
		return nil, err
	}

	toSave := dto.ConvertToExchangeRate(result)

	return toSave, s.repo.SaveRate(ctx, toSave)

}

func (s Service) HealthCheck(ctx context.Context) error {
	_, err := s.api.Call(ctx)

	if err != nil {
		s.logger.Error(
			"failed HealthCheck on API service call")
		return err
	}

	err = s.repo.CheckConnection(ctx)

	if err != nil {
		s.logger.Error("failed HealthCheck on connection to db")
	}
	return err
}
