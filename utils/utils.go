package utils

import "github.com/gin-gonic/gin"

func GetBaseUrl(c *gin.Context) (string, error) {
	// Получение базового URL из запроса
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	baseURL := scheme + "://" + c.Request.Host

	return baseURL, nil
}