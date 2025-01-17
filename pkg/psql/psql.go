package psql

import (
	"context"
	"database/sql"
	"log"

	"fmt"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
	"strings"
)

type Writer struct {
	logger *zap.Logger
}

func (w *Writer) Write(p []byte) (n int, err error) {
	msg := strings.TrimSpace(string(p))
	w.logger.Debug(msg)

	return len(p), nil
}

func Connect(ctx context.Context, logger *zap.Logger, opts ...OptionFunc) (*sql.DB, error) {
	options := &psqlOptions{}

	for _, opt := range opts {
		opt(options)
	}

	connString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		options.host,
		options.port,
		options.user,
		options.pass,
		options.database)

	logger.Debug(fmt.Sprintf("connString: %s", connString))

	db, err := sql.Open("postgres", connString)

	if err != nil {
		return nil, err
	}

	if err = db.PingContext(ctx); err != nil {
		return nil, err
	}

	log.SetOutput(&Writer{logger: logger})
	log.SetFlags(0)

	if err = goose.Up(db, options.migrations); err != nil {
		return nil, err
	}

	return db, nil

}
