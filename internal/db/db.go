package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type DataBaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Dbname   string
}

func (db *DataBaseConfig) LoadConfig() {
	// db.Host = "localhost"
	// db.Port = 5432
	// db.User = ";;;"
	// db.Password = ";;;"
	// db.Dbname = ";;;"
}

// func (db *DataBaseConfig) Close() error {
// 	return db.Close()
// }

func Connect(dbConf DataBaseConfig) (*sql.DB, error) {
	connectStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", dbConf.Host, dbConf.Port, dbConf.User, dbConf.Password, dbConf.Dbname)
	db, err := sql.Open("postgres", connectStr)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db, nil
}
