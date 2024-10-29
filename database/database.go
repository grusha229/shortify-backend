package database

import (
	"database/sql"
	"fmt"
	"log"
	"shortify/config"

	_ "github.com/lib/pq" // Подключение PostgreSQL драйвера
)

func Connect(cfg *config.Config) (*sql.DB, error) {
    connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

    db, err := sql.Open("postgres", connStr)
    if err != nil {
        return nil, err
    }
    err = db.Ping()
    if err != nil {
        return nil, err
    }
    createTable(db)
    return db, nil
}

func createTable(db *sql.DB) {
    query := `
    CREATE TABLE IF NOT EXISTS links (
        id SERIAL PRIMARY KEY,
        code VARCHAR(10) UNIQUE NOT NULL,
        url TEXT NOT NULL

    )`
    if _, err := db.Exec(query); err != nil {
        log.Fatal("Не удалось создать таблицу:", err)
    }
}
