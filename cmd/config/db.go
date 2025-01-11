package config

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

type Postgres struct {
	db *sql.DB
}

func NewPostgres() *Postgres {
	return &Postgres{}
}

func (p *Postgres) InitDB(cfg *Config) *Postgres {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("%v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("%v", err)
	}

	p.db = db

	return p
}

func (p *Postgres) GetDB() *sql.DB {
	return p.db
}

func (p *Postgres) Migrate() *Postgres {
	driver, err := postgres.WithInstance(p.db, &postgres.Config{})
	if err != nil {
		log.Fatalf("failed to create driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver,
	)

	if err != nil {
		log.Fatalf("failed to create migrate instance: %v", err)
	}

	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("failed to apply migrations: %v", err)
	}

	return p
}

func (p *Postgres) Close() {
	p.db.Close()
}
