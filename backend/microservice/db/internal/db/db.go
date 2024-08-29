package db

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Dbname   string
}

func Readfile(filePath string) (map[string]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("can't open file: %v", err)
	}
	defer file.Close()

	data := make(map[string]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			data[key] = value
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error read file: %v", err)
	}

	return data, nil
}

func Connect(filePath string) (*sql.DB, error) {
	fileData, err := Readfile(filePath)
	if err != nil {
		return nil, fmt.Errorf("cant't read file: %v", err)
	}

	port, err := strconv.Atoi(fileData["port"])
	if err != nil {
		return nil, fmt.Errorf("cant't read port: %v", err)
	}
	connectStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		fileData["host"], port, fileData["user"], fileData["password"], fileData["dbname"],
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
