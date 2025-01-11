package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/KirillGreenev/crypto-rate-service/internal/models"
)

type GarantexAPI interface {
	GetRates() (models.ResponseAPI, error)
}

type GarantexAPIImpl struct {
}

func NewGarantexApiImpl() *GarantexAPIImpl {
	return &GarantexAPIImpl{}
}

func (i *GarantexAPIImpl) GetRates() (models.ResponseAPI, error) {
	r, err := http.Get("https://garantex.org/api/v2/depth?market=usdtrub")
	if err != nil {
		return models.ResponseAPI{}, fmt.Errorf("repository.api.GetRates: %w", err)
	}
	defer r.Body.Close()

	var response models.ResponseAPI
	err = json.NewDecoder(r.Body).Decode(&response)
	if err != nil {
		return models.ResponseAPI{}, fmt.Errorf("repository.api.GetRates(): %w", err)
	}
	return response, nil
}
