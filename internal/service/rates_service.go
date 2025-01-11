package service

import (
	"context"
	"github.com/KirillGreenev/crypto-rate-service/internal/models"
	"github.com/KirillGreenev/crypto-rate-service/internal/repository"
	"github.com/KirillGreenev/crypto-rate-service/internal/repository/api"
	"time"
)

type RatesService interface {
	GetRates(ctx context.Context) (models.ResponseService, error)
}

type RatesServiceImpl struct {
	api  api.GarantexAPI
	repo repository.RatesRepository
}

func NewRatesServiceImpl(api api.GarantexAPI, repo repository.RatesRepository) *RatesServiceImpl {
	return &RatesServiceImpl{api: api, repo: repo}
}

func (r *RatesServiceImpl) GetRates(ctx context.Context) (models.ResponseService, error) {
	resultAPI, err := r.api.GetRates()
	if err != nil {
		return models.ResponseService{}, err
	}
	ask, bid := resultAPI.Asks[0], resultAPI.Bids[0]

	err = r.repo.Create(ctx, time.Unix(resultAPI.Timestamp, 0), ask, bid)
	if err != nil {
		return models.ResponseService{}, err
	}

	return models.ResponseService{resultAPI.Timestamp, ask, bid}, nil
}
