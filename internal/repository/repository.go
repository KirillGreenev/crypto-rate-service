package repository

import (
	"context"
	"github.com/KirillGreenev/crypto-rate-service/internal/models"
	"time"
)

type RatesRepository interface {
	Create(ctx context.Context, timestamp time.Time, ask models.Ask, bid models.Bid) error
}
