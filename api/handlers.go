package api

import (
	"database/sql"
	"net/http"
	"shortify/service"
	"shortify/utils"

	"github.com/gin-gonic/gin"
)

type ShortenRequest struct {
    URL string `json:"url"`
}

func CreateShortLink(c *gin.Context, db *sql.DB) {
    var request ShortenRequest
    println(&request)
    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "неправильный запрос"})
        return
    }

    code, err := service.CreateShortLink(db, request.URL)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось создать ссылку"})
        return
    }
    baseUrl, err := utils.GetBaseUrl(c);
    if err != nil {
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "short_url": baseUrl + "/link/" + code,
        "original_url": request.URL,
    })
}

func Redirect(c *gin.Context, db *sql.DB) {
    code := c.Param("code")
    url, err := service.GetOriginalURL(db, code)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "ссылка не найдена"})
        return
    }

    c.Redirect(http.StatusMovedPermanently, url)
}