package grpc

import (
	"go.uber.org/zap"
)

func New(opts ...Option) *Server {
	server := &Server{}

	for _, opt := range opts {
		opt(server)
	}

	return server
}

type Option func(*Server)

func WithLogger(logger *zap.Logger) Option {
	return func(s *Server) {
		s.logger = logger
	}
}

func WithExchangeRateService(service Service) Option {
	return func(s *Server) {
		s.service = service
	}
}
