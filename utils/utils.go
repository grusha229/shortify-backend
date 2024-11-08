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

func GetUTMParams(c *gin.Context) map[string]string {
    utmParams := []string{"utm_source", "utm_medium", "utm_campaign", "utm_term", "utm_content"}

    utmData := make(map[string]string)

    for _, param := range utmParams {
        value := c.Query(param)  // Получаем значение UTM-метки из запроса
        if value != "" {
            utmData[param] = value  // Сохраняем в карту, если параметр присутствует
        }
    }

    return utmData
}