package utils

import (
	"fmt"
	"log"
	"net"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/oschwald/geoip2-golang"
)

func GetBaseUrl(c *gin.Context) (string, error) {
	// Получаем значение заголовка Origin
	origin := c.Request.Header.Get("Origin")
	var scheme string

	if origin != "" {
		// Пытаемся извлечь схему из заголовка Origin
		parsedURL, err := url.Parse(origin)
		if err != nil {
			// Если не удалось разобрать URL, использует HTTP по умолчанию
			scheme = "http"
		} else {
			// Если Origin содержит схему, используем её
			scheme = parsedURL.Scheme
		}
	} else {
		// Если заголовка Origin нет, проверяем TLS соединение
		if c.Request.TLS != nil {
			scheme = "https"
		} else {
			scheme = "http"
		}
	}

	// Формируем базовый URL
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

func GetLocationByIp(ipAddress string) (string, error) {
    geoDB, err := geoip2.Open("GeoIP2-City.mmdb")

	if err != nil {
		log.Fatal(err)
	}
	defer geoDB.Close()
	ip := net.ParseIP(ipAddress)
    if ip == nil {
		return "", fmt.Errorf("invalid IP address: %s", ipAddress)
	}

	record, err := geoDB.City(ip)
	if err != nil {
		log.Fatal(err)
	}
    response := record.Country.Names["en"] + ", " + record.City.Names["en"]
    return response, nil
}