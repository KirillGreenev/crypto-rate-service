package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/KirillGreenev/crypto-rate-service/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockGarantexAPI struct {
	mock.Mock
}

func (m *MockGarantexAPI) GetRates() (models.ResponseAPI, error) {
	args := m.Called()
	return args.Get(0).(models.ResponseAPI), args.Error(1)
}

type MockRatesRepository struct {
	mock.Mock
}

func (m *MockRatesRepository) Create(ctx context.Context, timestamp time.Time, ask models.Ask, bid models.Bid) error {
	args := m.Called(ctx, timestamp, ask, bid)
	return args.Error(0)
}

func TestRatesServiceImpl_GetRates(t *testing.T) {
	mockAPI := new(MockGarantexAPI)
	mockRepo := new(MockRatesRepository)

	service := NewRatesServiceImpl(mockAPI, mockRepo)

	timestamp := time.Now().Unix()
	fakeResponse := models.ResponseAPI{
		Timestamp: timestamp,
		Asks: []models.Ask{
			{Price: "100.0", Volume: "1.0", Amount: "1.0", Factor: "1.0", Type: "ask"},
		},
		Bids: []models.Bid{
			{Price: "99.0", Volume: "2.0", Amount: "2.0", Factor: "1.0", Type: "bid"},
		},
	}

	mockAPI.On("GetRates").Return(fakeResponse, nil)
	mockRepo.On("Create", mock.Anything, time.Unix(timestamp, 0), fakeResponse.Asks[0], fakeResponse.Bids[0]).Return(nil)

	ctx := context.Background()
	response, err := service.GetRates(ctx)

	assert.NoError(t, err)
	assert.Equal(t, models.ResponseService{Timestamp: timestamp, Ask: fakeResponse.Asks[0], Bid: fakeResponse.Bids[0]}, response)

	mockAPI.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestRatesServiceImpl_GetRates_APIError(t *testing.T) {
	mockAPI := new(MockGarantexAPI)
	mockRepo := new(MockRatesRepository)

	service := NewRatesServiceImpl(mockAPI, mockRepo)

	mockAPI.On("GetRates").Return(models.ResponseAPI{}, errors.New("api error"))

	ctx := context.Background()
	response, err := service.GetRates(ctx)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "api error")
	assert.Equal(t, models.ResponseService{}, response)

	mockAPI.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestRatesServiceImpl_GetRates_RepoError(t *testing.T) {
	mockAPI := new(MockGarantexAPI)
	mockRepo := new(MockRatesRepository)

	service := NewRatesServiceImpl(mockAPI, mockRepo)

	timestamp := time.Now().Unix()
	fakeResponse := models.ResponseAPI{
		Timestamp: timestamp,
		Asks: []models.Ask{
			{Price: "100.0", Volume: "1.0", Amount: "1.0", Factor: "1.0", Type: "ask"},
		},
		Bids: []models.Bid{
			{Price: "99.0", Volume: "2.0", Amount: "2.0", Factor: "1.0", Type: "bid"},
		},
	}

	mockAPI.On("GetRates").Return(fakeResponse, nil)
	mockRepo.On("Create", mock.Anything, time.Unix(timestamp, 0), fakeResponse.Asks[0], fakeResponse.Bids[0]).Return(errors.New("repo error"))

	ctx := context.Background()
	response, err := service.GetRates(ctx)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "repo error")
	assert.Equal(t, models.ResponseService{}, response)

	mockAPI.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}
