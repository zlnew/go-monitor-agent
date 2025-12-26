package domain

import (
	"time"

	"github.com/google/uuid"
)

type EventLogReceived struct {
	ID        int64     `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Level     LogLevel  `json:"level"`
	Source    LogSource `json:"source"`
	Action    LogAction `json:"action"`
	TraceID   uuid.UUID `json:"trace_id"`

	ServerID      *uuid.UUID `json:"server_id"`
	ApplicationID *int64     `json:"application_id"`
	DeploymentID  *int64     `json:"deployment_id"`

	Message string      `json:"message"`
	Context *LogContext `json:"context"`
}
