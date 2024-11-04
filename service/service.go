package service

import (
	"database/sql"
	"errors"
	"sync"
)

// Кэш для хранения записей
var cache sync.Map

func CreateShortLink(db *sql.DB, originalURL string) (string, error) {
    // Проверка наличия записи в кэше
    if code, found := cache.Load(originalURL); found {
        return code.(string), nil
    }

    // Проверка наличия записи в базе данных
    var existingCode string
    err := db.QueryRow("SELECT code FROM links WHERE url = $1", originalURL).Scan(&existingCode)
    if err == nil {
        // Если запись найдена, добавляем её в кэш и возвращаем код
        cache.Store(originalURL, existingCode)
        return existingCode, nil
    } else if err != sql.ErrNoRows {
        // Если произошла ошибка при запросе, возвращаем её
        return "", err
    }

    // Если записи нет в базе, генерируем новый код
    code := GenerateCode(6)
    _, err = db.Exec("INSERT INTO links (code, url) VALUES ($1, $2)", code, originalURL)
    if err != nil {
        return "", err
    }

    // Сохраняем новый код в кэше
    cache.Store(originalURL, code)

    return code, nil
}

func GetOriginalURL(db *sql.DB, code string) (string, error) {
    var url string
    err := db.QueryRow("SELECT url FROM links WHERE code = $1", code).Scan(&url)
    if err != nil {
        if err == sql.ErrNoRows {
            return "", errors.New("ссылка не найдена")
        }
        return "", err
    }
    return url, nil
}
