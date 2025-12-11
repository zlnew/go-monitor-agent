package agent

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
)

type ServerCommand struct {
	Type    string          `json:"type"`
	Command string          `json:"command"`
	Payload json.RawMessage `json:"payload,omitempty"`
}

type DeploymentPayload struct {
	AppID   int64  `json:"app_id"`
	Version string `json:"version"`
	RepoURL string `json:"repo_url"`
}

func (a *Agent) handleCommand(ctx context.Context, cmd ServerCommand) {
	a.log.Info("executing command", "command", cmd.Command)

	switch cmd.Command {
	case "Ping":
		a.sendResponse(cmd.Command, map[string]string{"status": "pong"})

	case "DeployApplication":
		var payload DeploymentPayload
		if err := json.Unmarshal(cmd.Payload, &payload); err != nil {
			a.log.Error("invalid DeployApp payload", "error", err)
			a.sendResponse(cmd.Command, map[string]string{"status": "error", "message": "Invalid payload format"})
			return
		}
		a.log.Debug("received deployment payload", "app_id", payload.AppID, "version", payload.Version)
		a.sendResponse(cmd.Command, map[string]string{"status": "ok", "message": "Deployment command received"})

	case "StopApplication":
		a.log.Debug("received stop application command")
		a.sendResponse(cmd.Command, map[string]string{"status": "ok", "message": "Stop command received"})

	default:
		a.log.Warn("unknown command received", "command", cmd.Command)
		a.sendResponse(cmd.Command, map[string]string{"status": "error", "message": "Unknown command"})
	}
}

func (a *Agent) sendResponse(sourceCommand string, data any) {
	resp := struct {
		Type    string `json:"type"`
		Command string `json:"command"`
		Payload any    `json:"payload"`
	}{
		Type:    "response",
		Command: sourceCommand,
		Payload: data,
	}

	bytes, err := json.Marshal(resp)
	if err != nil {
		a.log.Error("failed to marshal command response", "error", err)
		return
	}

	a.conn.SetWriteDeadline(time.Now().Add(writeWait))
	if err := a.conn.WriteMessage(websocket.TextMessage, bytes); err != nil {
		a.log.Error("failed to write command response", "error", err)
	}
}
