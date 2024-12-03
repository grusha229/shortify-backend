package main

import (
	"log"
	"shortify/api"
	"shortify/config"
	"shortify/database"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
    // Загрузка конфигурации из .env
    cfg := config.LoadConfig()

    // Подключение к базе данных
    db, err := database.Connect(cfg)
    if err != nil {
        log.Fatal("Не удалось подключиться к базе данных:", err)
    }
    defer db.Close() // Закрытие подключения при завершении работы

    // Настройка маршрутов
    r := gin.Default()
    r.POST("api/shorten", func(c *gin.Context) { api.CreateShortLink(c, db) })
    r.GET("api/:code", func(c *gin.Context) { api.Redirect(c, db) })
    r.GET("api/details/:id", func(c *gin.Context) { api.GetLinkDetails(c, db) })

    // Запуск сервера
    if err := r.Run(":" + cfg.ServerPort); err != nil {
        log.Fatal("Ошибка при запуске сервера:", err)
    }
}