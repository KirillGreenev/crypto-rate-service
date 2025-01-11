package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/KirillGreenev/crypto-rate-service/internal/models"
)

type PostgesRepositoryImpl struct {
	db *sql.DB
}

func NewPostgesRepositoryImpl(db *sql.DB) *PostgesRepositoryImpl {
	return &PostgesRepositoryImpl{db: db}
}

func (p PostgesRepositoryImpl) Create(ctx context.Context, timestamp time.Time, ask models.Ask, bid models.Bid) error {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("repository.postgres.Create(): %w", err)
	}
	defer tx.Rollback()

	var askID, bidID int
	queryAsk := `INSERT INTO ask (price, volume, amount, factor, type) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err = tx.QueryRowContext(ctx, queryAsk,
		ask.Price, ask.Volume, ask.Amount, ask.Factor, ask.Type).Scan(&askID)
	if err != nil {
		return fmt.Errorf("repository.postgres.Create(): %w", err)
	}

	queryBid := `INSERT INTO bid (price, volume, amount, factor, type) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err = tx.QueryRowContext(ctx, queryBid,
		bid.Price, bid.Volume, bid.Amount, bid.Factor, bid.Type).Scan(&bidID)
	if err != nil {
		return fmt.Errorf("repository.postgres.Create(): %w", err)
	}

	queryRates := `INSERT INTO rates (timestamp, ask_id, bid_id) VALUES ($1, $2, $3)`
	_, err = tx.ExecContext(ctx, queryRates,
		timestamp, askID, bidID)
	if err != nil {
		return fmt.Errorf("repository.postgres.Create(): %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("repository.postgres.Create(): %w", err)
	}

	return nil
}
