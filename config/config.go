package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
    DBUser     string
    DBPassword string
    DBName     string
    DBHost     string
    DBPort     string
    ServerPort string
}

// LoadConfig загружает параметры конфигурации из переменных окружения
func LoadConfig() *Config {
    err := godotenv.Load("./dev.env");
    if err != nil {
        log.Fatal("Не удалось загрузить файл dev.env")
    }

    return &Config{
        DBUser:     os.Getenv("DB_USER"),
        DBPassword: os.Getenv("DB_PASSWORD"),
        DBName:     os.Getenv("DB_NAME"),
        DBHost:     os.Getenv("DB_HOST"),
        DBPort:     os.Getenv("DB_PORT"),
        ServerPort: os.Getenv("SERVER_PORT"),
    }
}
