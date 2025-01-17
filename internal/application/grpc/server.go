package grpc

import (
	"context"
	pb "github.com/peter-pashchenko/grpcExchangeService/internal/generated/grpc/exchangeRate"
	"github.com/peter-pashchenko/grpcExchangeService/internal/models"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"

	"google.golang.org/grpc/status"
)

//go:generate mockery --name=Service --case underscore  --with-expecter
type Service interface {
	GetRate(ctx context.Context) (*models.ExchangeRate, error)
	HealthCheck(ctx context.Context) error
}
type Server struct {
	pb.UnimplementedExchangeRateServiceServer

	service Service

	logger *zap.Logger
}

func (s *Server) GetRates(ctx context.Context, r *pb.EmptyRequest) (*pb.GetRatesResponse, error) {
	rate, err := s.service.GetRate(ctx)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot get rate")
	}

	return &pb.GetRatesResponse{
		Timestamp: rate.Timestamp,
		AskPrice:  rate.AskPrice,
		BidPrice:  rate.BidPrice}, nil
}

func (s *Server) HealthCheck(ctx context.Context, r *pb.EmptyRequest) (*pb.HealthCheckResponse, error) {
	err := s.service.HealthCheck(ctx)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "healthcheck didn't pass")
	}

	return &pb.HealthCheckResponse{Healthcheckpass: true}, nil
}
