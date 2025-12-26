package domain

import (
	"encoding/json"

	"github.com/google/uuid"
)

type WsClientMessage struct {
	Type    string          `json:"type"`
	Channel string          `json:"channel,omitempty"`
	Event   string          `json:"event,omitempty"`
	Payload json.RawMessage `json:"payload,omitempty"`
}

type WsServerEvent struct {
	Channel string `json:"channel"`
	Event   string `json:"event"`
	Payload any    `json:"payload,omitempty"`
}

type WsServerMessage struct {
	TargetServerID uuid.UUID       `json:"target_server_id"`
	Payload        json.RawMessage `json:"payload"`
}
