package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"

	"amocrm-tg-bot/internal/config"
	"amocrm-tg-bot/internal/handler"
	"amocrm-tg-bot/internal/tg"
)

func main() {
	// Загрузка переменных окружения из .env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	cfg := config.Load()

	// Инициализация Telegram-бота
	sender := tg.New(cfg.TelegramToken)
	go sender.StartPolling() // запуск обработчика команд Telegram

	// HTTP обработчик для /notify
	h := handler.New(sender, cfg.ManagerChatID)
	http.HandleFunc("/notify", h.Notify)

	log.Printf("Starting server on port %s...", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, nil); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
