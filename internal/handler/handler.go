package handler

import (
	"amocrm-tg-bot/internal/parser"
	"amocrm-tg-bot/internal/tg"
	"fmt"
	"io"
	"net/http"
)

type Handler struct {
	Sender        *tg.Sender
	ManagerChatID string
}

func New(sender *tg.Sender, chatID string) *Handler {
	return &Handler{
		Sender:        sender,
		ManagerChatID: chatID,
	}
}

func (h *Handler) Notify(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	event, err := parser.ParseEvent(body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Bad event: %v", err), http.StatusBadRequest)
		return
	}

	chatID, err := tg.LoadManagerID()
	if err != nil {
		http.Error(w, fmt.Sprintf("Cannot load manager ID: %v", err), http.StatusInternalServerError)
		return
	}

	if err := h.Sender.SendMessage(chatID, event.Text); err != nil {
		http.Error(w, fmt.Sprintf("Failed to send message: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func (h *Handler) Start(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
		return
	}

	chatID, err := tg.LoadManagerID()
	if err != nil {
		http.Error(w, fmt.Sprintf("Cannot load manager ID: %v", err), http.StatusInternalServerError)
		return
	}

	msg := "🤖 Бот готов принимать события от AmoCRM!"
	if err := h.Sender.SendMessage(chatID, msg); err != nil {
		http.Error(w, fmt.Sprintf("Error sending start message: %v", err), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Start acknowledged"))
}
