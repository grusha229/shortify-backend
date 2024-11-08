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
    // Создание расширения pgcrypto для генерации UUID
    _, err := db.Exec("CREATE EXTENSION IF NOT EXISTS pgcrypto;")
    if err != nil {
        log.Fatal("Не удалось создать расширение pgcrypto:", err)
    }

    // Создание таблицы links с UUID для поля id
    linksQuery := `
    CREATE TABLE IF NOT EXISTS links (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        code VARCHAR(10) UNIQUE NOT NULL,
        url TEXT NOT NULL
    )`

    if _, err := db.Exec(linksQuery); err != nil {
        log.Fatal("Не удалось создать таблицу links:", err)
    }

    // Создание таблицы link_visits с UUID для поля link_id
    statQuery := `
    CREATE TABLE IF NOT EXISTS link_visits (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        link_id UUID REFERENCES links(id) ON DELETE CASCADE,  -- Используем UUID для связи с links
        visited_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        ip_address VARCHAR(45),
        user_agent TEXT,
        utm_source VARCHAR(255)  -- Поле для хранения UTM-метки
    )`

    if _, err := db.Exec(statQuery); err != nil {
        log.Fatal("Не удалось создать таблицу link_visits:", err)
    }
}
