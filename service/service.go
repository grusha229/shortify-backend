package service

import (
	"database/sql"
	"errors"
)

func CreateShortLink(db *sql.DB, originalURL string) (string, error) {
    code := GenerateCode(6)
    _, err := db.Exec("INSERT INTO links (code, url) VALUES ($1, $2)", code, originalURL)
    if err != nil {
        return "", err
    }
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
