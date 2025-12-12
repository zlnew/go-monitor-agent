package agent

import (
	"strconv"

	"horizonx-server/internal/domain"
)

func (a *Agent) handleCommand(cmd domain.WsAgentCommand) {
	a.log.Info("executing command", "command", cmd.CommandType)

	switch cmd.CommandType {
	case "init":
		a.initializeAgent(cmd.TargetServerID)
	default:
		a.log.Warn("unknown command received", "command", cmd.CommandType)
	}
}

func (a *Agent) initializeAgent(serverID string) {
	id, err := strconv.ParseInt(serverID, 10, 64)
	if err != nil {
		a.log.Error("failed to parse server id to int64")
	}

	select {
	case a.initCh <- id:
		a.log.Debug("received server id via init command")
	default:
		a.log.Warn("init channel full, dropping init ID (should not happen)")
	}
}
