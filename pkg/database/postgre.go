package database

import (
	"SIGE/pkg/config"

	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq" // <- важен символ "_"!
)

var DB *sql.DB

func InitDB() {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.GetEnv("DB_HOST"),
		config.GetEnv("DB_PORT"),
		config.GetEnv("DB_USER"),
		config.GetEnv("DB_PASSWORD"),
		config.GetEnv("DB_NAME"),
	)
	var err error
	DB, err = sql.Open("postgres", dsn) // присваиваем глобальной переменной
	if err != nil {
		log.Fatal(err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully connected to PostgresSQL!")
}
