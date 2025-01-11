package postgres

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/KirillGreenev/crypto-rate-service/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestPostgesRepositoryImpl_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock database: %v", err)
	}
	defer db.Close()

	repo := NewPostgesRepositoryImpl(db)

	timestamp := time.Now()
	ask := models.Ask{Price: "100.0", Volume: "1.0", Amount: "1.0", Factor: "1.0", Type: "ask"}
	bid := models.Bid{Price: "99.0", Volume: "2.0", Amount: "2.0", Factor: "1.0", Type: "bid"}

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO ask \(price, volume, amount, factor, type\) VALUES \(\$1, \$2, \$3, \$4, \$5\) RETURNING id`).
		WithArgs(ask.Price, ask.Volume, ask.Amount, ask.Factor, ask.Type).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	mock.ExpectQuery(`INSERT INTO bid \(price, volume, amount, factor, type\) VALUES \(\$1, \$2, \$3, \$4, \$5\) RETURNING id`).
		WithArgs(bid.Price, bid.Volume, bid.Amount, bid.Factor, bid.Type).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(2))

	mock.ExpectExec(`INSERT INTO rates \(timestamp, ask_id, bid_id\) VALUES \(\$1, \$2, \$3\)`).
		WithArgs(timestamp, 1, 2).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectCommit()

	ctx := context.Background()
	err = repo.Create(ctx, timestamp, ask, bid)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPostgesRepositoryImpl_Create_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock database: %v", err)
	}
	defer db.Close()

	repo := NewPostgesRepositoryImpl(db)

	timestamp := time.Now()
	ask := models.Ask{Price: "100.0", Volume: "1.0", Amount: "1.0", Factor: "1.0", Type: "ask"}
	bid := models.Bid{Price: "99.0", Volume: "2.0", Amount: "2.0", Factor: "1.0", Type: "bid"}

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO ask \(price, volume, amount, factor, type\) VALUES \(\$1, \$2, \$3, \$4, \$5\) RETURNING id`).
		WithArgs(ask.Price, ask.Volume, ask.Amount, ask.Factor, ask.Type).
		WillReturnError(sql.ErrNoRows) // Эмулируем ошибку

	ctx := context.Background()
	err = repo.Create(ctx, timestamp, ask, bid)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "repository.postgres.Create():")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
