package api

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	LinksModels "shortify/models"
	"shortify/service"
	"shortify/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateShortLink(c *gin.Context, db *sql.DB) {
    var request LinksModels.ShortenRequest
    println(&request)
    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "неправильный запрос"})
        return
    }

    data, err := service.CreateShortLink(db, request.URL)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось создать ссылку", "text": err.Error()})
        return
    }
    baseUrl, err := utils.GetBaseUrl(c);
    if err != nil {
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "id": data.ID,
        "short_url": baseUrl + "/link/" + data.ShortURL,
        "original_url": data.OriginalURL,
    })
}

func GetLinkDetails(c *gin.Context, db *sql.DB) {
    LinkId := c.Param("id")

    parsedUUID, err := uuid.Parse(LinkId)
    if err != nil {
        log.Printf("Invalid UUID: %v\n", err)
        c.JSON(400, gin.H{"error": "Invalid UUID format"})
        return
    }

    data, err := service.GetURLStatData(db, parsedUUID)

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось получить данные по url", "text": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"data": data})
}

func Redirect(c *gin.Context, db *sql.DB) {
    code := c.Param("code")
    data, err := service.GetURLData(db, code)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "ссылка не найдена"})
        return
    }

    ipAddress := c.ClientIP()
    userAgent := c.GetHeader("User-Agent")
    utmParams := utils.GetUTMParams(c)
    utmSource := utmParams["utm_source"]
    fmt.Printf("Visit recoded")
    err = service.RecordVisit(db, data.ID, ipAddress, userAgent, utmSource)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err})
        return
    }

    c.Header("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
    c.Header("Pragma", "no-cache")
    c.Header("Expires", "0")

    c.Redirect(http.StatusMovedPermanently, data.OriginalURL)
}