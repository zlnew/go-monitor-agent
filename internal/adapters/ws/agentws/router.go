package agentws

import (
	"context"
	"encoding/json"

	"horizonx-server/internal/domain"
	"horizonx-server/internal/logger"

	"github.com/google/uuid"
)

type Router struct {
	ctx    context.Context
	cancel context.CancelFunc

	agents map[uuid.UUID]*Client

	register   chan *Client
	unregister chan *Client
	commands   chan *domain.WsAgentCommand

	log logger.Logger
}

func NewRouter(parent context.Context, log logger.Logger) *Router {
	ctx, cancel := context.WithCancel(parent)

	return &Router{
		ctx:        ctx,
		cancel:     cancel,
		agents:     make(map[uuid.UUID]*Client),
		register:   make(chan *Client, 64),
		unregister: make(chan *Client, 64),
		commands:   make(chan *domain.WsAgentCommand, 1024),
		log:        log,
	}
}

func (r *Router) Run() {
	for {
		select {
		case <-r.ctx.Done():
			r.log.Info("ws: agent hub shutting down...")
			for _, agent := range r.agents {
				close(agent.send)
			}
			return

		case a := <-r.register:
			r.agents[a.ID] = a
			a.log.Info("ws: agent registered", "id", a.ID)

		case a := <-r.unregister:
			agent, ok := r.agents[a.ID]
			if !ok {
				continue
			}

			delete(r.agents, a.ID)
			close(agent.send)
			r.log.Info("ws: agent unregistered", "id", a.ID)

		case cmd := <-r.commands:
			r.handleCommand(cmd)
		}
	}
}

func (r *Router) Stop() {
	r.cancel()
}

func (r *Router) SendCommand(cmd *domain.WsAgentCommand) {
	select {
	case r.commands <- cmd:
	case <-r.ctx.Done():
	default:
		r.log.Warn("ws: command buffer full, dropping command", "command", cmd.CommandType)
	}
}

func (r *Router) handleCommand(cmd *domain.WsAgentCommand) {
	agent, ok := r.agents[cmd.TargetServerID]
	if !ok {
		r.log.Warn("ws: target agent not connected", "server_id", cmd.TargetServerID)
		return
	}

	message, err := json.Marshal(cmd)
	if err != nil {
		r.log.Error("ws: failed to marshal server command", "error", err)
		return
	}

	select {
	case agent.send <- message:
		r.log.Info("ws: command sent to agent", "server_id", agent.ID, "command", cmd.CommandType)
	default:
		r.log.Warn("ws: agent channel full, force unregister", "server_id", agent.ID)
		r.unregister <- agent
	}
}
