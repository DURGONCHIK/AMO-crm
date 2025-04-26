package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	TelegramToken string
	ManagerChatID string
	Port          string
}

func Load() *Config {
	// Загружаем .env файл
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("⚠️ .env file not found, reading from environment variables")
	}

	// Чтение переменных окружения
	token := os.Getenv("TG_BOT_TOKEN")
	if token == "" {
		log.Fatal("❌ TG_BOT_TOKEN is required but not set")
	}

	chatID := os.Getenv("MANAGER_CHAT_ID")
	if chatID == "" {
		log.Fatal("❌ MANAGER_CHAT_ID is required but not set")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return &Config{
		TelegramToken: token,
		ManagerChatID: chatID,
		Port:          port,
	}
}
