package service

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	LinksModels "shortify/models"
	"sync"

	"github.com/google/uuid"
)

// Кэш для хранения записей
var cache sync.Map

func CreateShortLink(db *sql.DB, originalURL string) (LinksModels.ShortenResponse, error) {
    if cachedData, found := cache.Load(originalURL); found {
        return cachedData.(LinksModels.ShortenResponse), nil
    }

    var existingCode string
    var existingID uuid.UUID

    err := db.QueryRow("SELECT id, code FROM links WHERE url = $1", originalURL).Scan(&existingID, &existingCode)

    if err == nil {
        // Если запись найдена, добавляем её в кэш и возвращаем данные
        response := LinksModels.ShortenResponse{
            ID:          existingID,
            ShortURL:    existingCode,
            OriginalURL: originalURL,
        }
        cache.Store(originalURL, response)
    } else if err != sql.ErrNoRows {
        fmt.Printf("Error query: %v\n", err)
        return LinksModels.ShortenResponse{}, err
    }

    // Если записи нет в базе, генерируем новый код
    code := GenerateCode(6)
    newID := uuid.New()

    _, err = db.Exec("INSERT INTO links (id, code, url) VALUES ($1, $2, $3)", newID, code, originalURL)
    if err != nil {
        fmt.Printf("Error insert: %v\n", err)
        return LinksModels.ShortenResponse{}, err
    }

    response := LinksModels.ShortenResponse{
        ID:          newID,
        ShortURL:    code,
        OriginalURL: originalURL,
    }

    // Сохраняем новый код в кэше
    cache.Store(originalURL, response)

    return response, nil
}

func GetURLData(db *sql.DB, code string) (LinksModels.ShortenResponse, error) {
    var url string
    var id uuid.UUID
    err := db.QueryRow("SELECT url, id FROM links WHERE code = $1", code).Scan(&url, &id)
    if err != nil {
        if err == sql.ErrNoRows {
            return LinksModels.ShortenResponse{}, errors.New("ссылка не найдена")
        }
        return LinksModels.ShortenResponse{}, err
    }

    response := LinksModels.ShortenResponse{
        ID:          id,
        ShortURL:    code,
        OriginalURL: url,
    }
    return response, nil
}

func RecordVisit(db *sql.DB, linkID uuid.UUID, ipAddress string, userAgent string, utmSource string) error {
    fmt.Println("Type of linkID:", reflect.TypeOf(linkID))
    fmt.Println("Type of ipAddress:", reflect.TypeOf(ipAddress))
    fmt.Println("Type of userAgent:", reflect.TypeOf(userAgent))
    fmt.Println("Type of utmSource:", reflect.TypeOf(utmSource))

    _, err := db.Exec(`
        INSERT INTO link_visits (link_id, ip_address, user_agent, utm_source)
        VALUES ($1, $2, $3, $4)
    `, linkID, ipAddress, userAgent, utmSource)
    return err
}
