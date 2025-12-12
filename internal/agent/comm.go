package agent

import (
	"context"
	"encoding/json"
	"time"

	"horizonx-server/internal/domain"

	"github.com/gorilla/websocket"
)

func (a *Agent) readPump(ctx context.Context) error {
	a.conn.SetReadLimit(maxMessageSize)
	a.conn.SetReadDeadline(time.Now().Add(pongWait))
	a.conn.SetPongHandler(func(string) error {
		a.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		default:
			_, message, err := a.conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					a.log.Error("ws read error (unexpected close)", "error", err)
					return err
				} else {
					a.log.Info("ws read finished (normal closure or ping/pong timeout)")
					return nil
				}
			}

			var command domain.WsAgentCommand
			if err := json.Unmarshal(message, &command); err != nil {
				a.log.Error("invalid command payload received", "error", err)
				continue
			}

			a.handleCommand(command)
		}
	}
}

func (a *Agent) writePump(ctx context.Context) error {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case metrics, ok := <-a.metricsCh:
			if !ok {
				a.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return nil
			}

			a.conn.SetWriteDeadline(time.Now().Add(writeWait))

			channel := domain.GetServerMetricsChannel(metrics.ServerID)
			event := domain.WsEventServerMetricsReport
			payload, err := json.Marshal(metrics)
			if err != nil {
				a.log.Error("failed to marshal metrics", "error", err)
				continue
			}

			message := &domain.WsClientMessage{
				Type:    domain.WsAgentReport,
				Channel: channel,
				Event:   event,
				Payload: payload,
			}

			bytes, err := json.Marshal(message)
			if err != nil {
				a.log.Error("failed to marshal metrics report message", "error", err)
				continue
			}

			if err := a.conn.WriteMessage(websocket.TextMessage, bytes); err != nil {
				a.log.Error("failed to write metrics report", "error", err)
				return err
			}

		case <-ticker.C:
			a.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := a.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				a.log.Error("failed to write ping", "error", err)
				return err
			}
		}
	}
}
