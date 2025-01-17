package exchangeRatesService

import (
	"context"
	"errors"
	"github.com/peter-pashchenko/grpcExchangeService/internal/application/dto"
	"github.com/peter-pashchenko/grpcExchangeService/internal/models"
	"github.com/peter-pashchenko/grpcExchangeService/internal/modules/services/exchangeRate/mocks"
	"github.com/peter-pashchenko/grpcExchangeService/pkg/logger"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_GetRate(t *testing.T) {
	logg := logger.New("debug")
	testCtx := context.TODO()

	dtoResults := dto.ExchangeRatesDTO{
		Timestamp: 12345,
		Asks:      []dto.Ask{{Price: "101"}, {Price: "2"}, {Price: "3"}},
		Bids:      []dto.Bid{{Price: "111"}, {Price: "1"}, {Price: "1"}},
	}

	modelsResults := models.ExchangeRate{
		Timestamp: 12345,
		AskPrice:  "101",
		BidPrice:  "111",
	}

	tests := []struct {
		name        string
		mockAPI     func() *mocks.Api
		mockRepo    func() *mocks.Repo
		expected    *models.ExchangeRate
		expectedERR bool
	}{
		{
			name: "success",
			mockAPI: func() *mocks.Api {
				mock := &mocks.Api{}
				mock.EXPECT().Call(testCtx).Return(&dtoResults, nil)
				return mock
			},
			mockRepo: func() *mocks.Repo {
				mock := &mocks.Repo{}
				mock.EXPECT().SaveRate(testCtx, &modelsResults).Return(nil)
				return mock
			},
			expected:    &models.ExchangeRate{Timestamp: 12345, AskPrice: "101", BidPrice: "111"},
			expectedERR: false,
		},
		{
			name: "error in api call",
			mockAPI: func() *mocks.Api {
				mock := &mocks.Api{}
				mock.EXPECT().Call(testCtx).Return(nil, errors.New("some error"))
				return mock
			},
			mockRepo:    func() *mocks.Repo { return nil },
			expected:    nil,
			expectedERR: true,
		},
		{
			name: "error in db call",
			mockAPI: func() *mocks.Api {
				mock := &mocks.Api{}
				mock.EXPECT().Call(testCtx).Return(&dtoResults, nil)
				return mock
			},
			mockRepo: func() *mocks.Repo {
				mock := &mocks.Repo{}
				mock.EXPECT().SaveRate(testCtx, &modelsResults).Return(errors.New("some error"))
				return mock
			},
			expected:    nil,
			expectedERR: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockService := New(test.mockRepo(), test.mockAPI(), logg)

			result, err := mockService.GetRate(testCtx)

			if test.expectedERR {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expected, result)
			}

		})
	}

}

func Test_HealthCheck(t *testing.T) {
	logg := logger.New("debug")
	testCtx := context.TODO()

	tests := []struct {
		name        string
		mockAPI     func() *mocks.Api
		mockRepo    func() *mocks.Repo
		expectedERR bool
	}{
		{
			name: "success",
			mockAPI: func() *mocks.Api {
				mock := &mocks.Api{}
				mock.EXPECT().Call(testCtx).Return(nil, nil)
				return mock
			},
			mockRepo: func() *mocks.Repo {
				mock := &mocks.Repo{}
				mock.EXPECT().CheckConnection(testCtx).Return(nil)
				return mock
			},
			expectedERR: false,
		},
		{
			name: "error in api call",
			mockAPI: func() *mocks.Api {
				mock := &mocks.Api{}
				mock.EXPECT().Call(testCtx).Return(nil, errors.New("some error"))
				return mock
			},
			mockRepo:    func() *mocks.Repo { return nil },
			expectedERR: true,
		},
		{
			name: "error in db call",
			mockAPI: func() *mocks.Api {
				mock := &mocks.Api{}
				mock.EXPECT().Call(testCtx).Return(nil, nil)
				return mock
			},
			mockRepo: func() *mocks.Repo {
				mock := &mocks.Repo{}
				mock.EXPECT().CheckConnection(testCtx).Return(errors.New("some error"))
				return mock
			},
			expectedERR: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockService := New(test.mockRepo(), test.mockAPI(), logg)

			err := mockService.HealthCheck(testCtx)

			if test.expectedERR {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

		})
	}

}
