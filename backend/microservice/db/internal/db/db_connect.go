package db

import (
	"InterestingChats/backend/microservice/db/internal/config"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func Connect(config config.DatabaseConfig) (*sql.DB, error) {
	connectStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.Dbname,
	)
	db, err := sql.Open("postgres", connectStr)
	if err != nil {
		return nil, fmt.Errorf("can't open database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("can't ping database: %v", err)
	}
	return db, nil
}
