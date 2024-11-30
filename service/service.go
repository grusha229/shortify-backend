package service

import (
	"database/sql"
	"errors"
	"fmt"
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

func GetURLStatData(db *sql.DB, id uuid.UUID) ([]LinksModels.GetURLStatDataResponse, error) {
    fmt.Println("Query ID:", id)

    query := `SELECT * FROM link_visits WHERE link_id = $1`
    rows, err := db.Query(query, id)
    fmt.Println("Query:", query)
    
    if err != nil {
        fmt.Println("Query error:", err)
        return nil, err
    }
    defer rows.Close()

    fmt.Println("Rows:", rows)

    var statistics []LinksModels.GetURLStatDataResponse

    for rows.Next() {
        var statDBData LinksModels.GetURLStatDataDB
        
        // Правильное сканирование данных
        if err := rows.Scan(&statDBData.ID, &statDBData.LinkId, &statDBData.VisitedAt, &statDBData.IPAddress, &statDBData.UserAgent, &statDBData.UtmSource); err != nil {
            fmt.Println("Scan error:", err)
            return nil, err
        }
        
        fmt.Printf("Scanned data: %+v\n", statDBData)

        // location, err := utils.GetLocationByIp(statDBData.IPAddress)
        if err != nil {
            fmt.Println("Location error:", err)
            return nil, err
        }

        fmt.Printf("Location for IP %s: %s\n", statDBData.IPAddress, statDBData.IPAddress)

        statistics = append(statistics, LinksModels.GetURLStatDataResponse{
            ID: statDBData.ID,
            Location: statDBData.IPAddress,
            LinkId: statDBData.LinkId,
            VisitedAt: statDBData.VisitedAt,
            UtmSource: statDBData.UtmSource,
        })
    }

    if err := rows.Err(); err != nil {
        fmt.Println("Rows iteration error:", err)
        return nil, err
    }

    fmt.Printf("Final statistics data: %+v\n", statistics)
    return statistics, nil
}


func RecordVisit(db *sql.DB, linkID uuid.UUID, ipAddress string, userAgent string, utmSource string) error {
    _, err := db.Exec(`
        INSERT INTO link_visits (link_id, ip_address, user_agent, utm_source)
        VALUES ($1, $2, $3, $4)
    `, linkID, ipAddress, userAgent, utmSource)
    return err
}
