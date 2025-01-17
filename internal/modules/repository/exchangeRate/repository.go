package exchangeRate

import (
	"context"
	"database/sql"
	"github.com/peter-pashchenko/grpcExchangeService/internal/models"
	"go.uber.org/zap"
)

type Repo struct {
	db *sql.DB

	logger *zap.Logger
}

func New(db *sql.DB, logger *zap.Logger) *Repo {
	return &Repo{db: db, logger: logger}
}

func (r *Repo) SaveRate(ctx context.Context, rate *models.ExchangeRate) error {
	query := `INSERT INTO exchangeRate (requestTimeStamp,askPrice,bidPrice) VALUES ($1,$2,$3)`

	r.logger.Debug("query string is ready", zap.String("query", query))

	_, err := r.db.ExecContext(
		ctx,
		query,
		rate.Timestamp,
		rate.AskPrice,
		rate.BidPrice)

	if err == nil {
		r.logger.Debug("data saved to db successfully")
	}

	return err
}

func (r *Repo) CheckConnection(ctx context.Context) error {
	return r.db.PingContext(ctx)
}
