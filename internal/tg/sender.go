package tg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Sender struct {
	Token string
}

func New(token string) *Sender {
	return &Sender{Token: token}
}

func (s *Sender) SendMessage(chatID string, text string) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", s.Token)

	body, _ := json.Marshal(map[string]string{
		"chat_id": chatID,
		"text":    text,
	})

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("telegram returned status %d", resp.StatusCode)
	}

	return nil
}

func (s *Sender) StartPolling() {
	offset := 0

	for {
		url := fmt.Sprintf("https://api.telegram.org/bot%s/getUpdates?timeout=60&offset=%d", s.Token, offset)
		resp, err := http.Get(url)
		if err != nil {
			log.Printf("Failed to get updates: %v", err)
			continue
		}

		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		var result struct {
			OK     bool `json:"ok"`
			Result []struct {
				UpdateID int `json:"update_id"`
				Message  struct {
					MessageID int `json:"message_id"`
					Chat      struct {
						ID int64 `json:"id"`
					} `json:"chat"`
					Text string `json:"text"`
				} `json:"message"`
			} `json:"result"`
		}

		if err := json.Unmarshal(body, &result); err != nil {
			log.Printf("Failed to parse response: %v", err)
			continue
		}

		for _, update := range result.Result {
			offset = update.UpdateID + 1

			if update.Message.Text == "/start" {
				managerID := update.Message.Chat.ID
				log.Printf("New message from chat ID: %d", managerID)

				err := ensureManagerIDFile(managerID)
				if err != nil {
					log.Printf("Failed to save manager ID: %v", err)
				}

				err = s.SendMessage(strconv.FormatInt(managerID, 10), "🤖 Привет! Бот готов работать!")
				if err != nil {
					log.Printf("Failed to send message: %v", err)
				}
			}
		}
	}
}

func ensureManagerIDFile(chatID int64) error {
	const filename = "manager_id.txt"

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return os.WriteFile(filename, []byte(strconv.FormatInt(chatID, 10)), 0644)
	}
	return nil
}

func LoadManagerID() (string, error) {
	data, err := os.ReadFile("manager_id.txt")
	if err != nil {
		return "", fmt.Errorf("failed to read manager_id.txt: %w", err)
	}
	return string(data), nil
}
