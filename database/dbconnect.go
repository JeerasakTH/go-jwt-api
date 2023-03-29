package database

import (
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

type PostgresDB struct {
	DB *sqlx.DB
}

func ConnectDB() (*PostgresDB, error) {
	_ = godotenv.Load()
	PostgresURL := os.Getenv("POSTGRES_URI")
	db, err := sqlx.Open("pgx", PostgresURL)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &PostgresDB{DB: db}, nil
}

func (p *PostgresDB) Close() {
	p.DB.Close()
}
