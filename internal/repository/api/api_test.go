package api

import (
	"testing"

	"github.com/KirillGreenev/crypto-rate-service/internal/models"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestGarantexAPIImpl_GetRates(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	t.Run("success", func(t *testing.T) {
		fakeResponse := models.ResponseAPI{
			Asks: []models.Ask{
				{Price: "100.0", Amount: "1.0"},
			},
			Bids: []models.Bid{
				{Price: "99.0", Amount: "2.0"},
			},
		}

		httpmock.RegisterResponder("GET", "https://garantex.org/api/v2/depth?market=usdtrub",
			httpmock.NewJsonResponderOrPanic(200, fakeResponse))

		api := NewGarantexApiImpl()
		response, err := api.GetRates()

		assert.NoError(t, err)
		assert.Equal(t, fakeResponse, response)
	})

	t.Run("http error", func(t *testing.T) {
		httpmock.RegisterResponder("GET", "https://garantex.org/api/v2/depth?market=usdtrub",
			httpmock.NewStringResponder(500, ""))

		api := NewGarantexApiImpl()
		_, err := api.GetRates()

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "repository.api.GetRates():")
	})

	t.Run("json decoding error", func(t *testing.T) {
		httpmock.RegisterResponder("GET", "https://garantex.org/api/v2/depth?market=usdtrub",
			httpmock.NewStringResponder(200, "invalid json"))

		api := NewGarantexApiImpl()
		_, err := api.GetRates()

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "repository.api.GetRates():")
	})
}
