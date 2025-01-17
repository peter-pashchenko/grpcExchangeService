package grpc

import (
	"context"
	"errors"
	"github.com/peter-pashchenko/grpcExchangeService/internal/application/grpc/mocks"
	pb "github.com/peter-pashchenko/grpcExchangeService/internal/generated/grpc/exchangeRate"
	"github.com/peter-pashchenko/grpcExchangeService/internal/models"
	"github.com/peter-pashchenko/grpcExchangeService/pkg/logger"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_GetRates(t *testing.T) {
	testCTX := context.TODO()
	logg := logger.New("debug")

	rate := models.ExchangeRate{
		Timestamp: 12345,
		AskPrice:  "100",
		BidPrice:  "101",
	}

	tests := []struct {
		name        string
		mockService func() *mocks.Service
		expected    *pb.GetRatesResponse
		expectedErr bool
	}{
		{
			name: "Success",
			mockService: func() *mocks.Service {
				mock := &mocks.Service{}
				mock.EXPECT().GetRate(testCTX).Return(&rate, nil)
				return mock
			},
			expected: &pb.GetRatesResponse{
				Timestamp: rate.Timestamp,
				AskPrice:  rate.AskPrice,
				BidPrice:  rate.BidPrice,
			},
			expectedErr: false,
		},
		{
			name: "Fail",
			mockService: func() *mocks.Service {
				mock := &mocks.Service{}
				mock.EXPECT().GetRate(testCTX).Return(nil, errors.New("some error"))
				return mock
			},
			expected:    nil,
			expectedErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			service := New(
				WithLogger(logg),
				WithExchangeRateService(test.mockService()))

			res, err := service.GetRates(testCTX, &pb.EmptyRequest{})

			if test.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expected, res)
			}
		})
	}

}

func Test_HealthCheck(t *testing.T) {
	testCTX := context.TODO()
	logg := logger.New("debug")

	tests := []struct {
		name        string
		mockService func() *mocks.Service
		expected    *pb.HealthCheckResponse
		expectedErr bool
	}{
		{
			name: "Success",
			mockService: func() *mocks.Service {
				mock := &mocks.Service{}
				mock.EXPECT().HealthCheck(testCTX).Return(nil)
				return mock
			},
			expected:    &pb.HealthCheckResponse{Healthcheckpass: true},
			expectedErr: false,
		},
		{
			name: "Fail",
			mockService: func() *mocks.Service {
				mock := &mocks.Service{}
				mock.EXPECT().HealthCheck(testCTX).Return(errors.New("some error"))
				return mock
			},
			expected:    nil,
			expectedErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			service := New(
				WithLogger(logg),
				WithExchangeRateService(test.mockService()))

			res, err := service.HealthCheck(testCTX, &pb.EmptyRequest{})

			if test.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expected, res)
			}
		})
	}

}
