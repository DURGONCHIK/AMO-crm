package parser

import (
	"encoding/json"
	"errors"
	"strings"
)

type EventType int

const (
	EventUnknown EventType = iota
	EventMessageNotDelivered
	EventClientOrder
)

type Event struct {
	Type EventType
	Text string
}

type AmoPayload struct {
	Event  string            `json:"event"`
	Fields map[string]string `json:"fields"`
}

func ParseEvent(data []byte) (*Event, error) {
	var payload AmoPayload
	if err := json.Unmarshal(data, &payload); err != nil {
		return nil, errors.New("invalid JSON payload")
	}

	switch payload.Event {
	case "message_not_delivered":
		return &Event{
			Type: EventMessageNotDelivered,
			Text: "❗ Сообщение клиенту в WhatsApp не доставлено. Свяжись с ним.",
		}, nil

	case "client_order":
		var builder strings.Builder
		builder.WriteString("🛒 Новый заказ от клиента:\n")

		for key, value := range payload.Fields {
			if strings.TrimSpace(value) != "" {
				builder.WriteString("- " + key + ": " + value + "\n")
			}
		}

		return &Event{
			Type: EventClientOrder,
			Text: builder.String(),
		}, nil
	}

	return nil, errors.New("unsupported event type")
}
