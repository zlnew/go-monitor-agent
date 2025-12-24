package domain

import "github.com/google/uuid"

type EventJobCreated struct {
	JobID         int64     `json:"job_id"`
	ServerID      uuid.UUID `json:"server_id"`
	ApplicationID *int64    `json:"application_id"`
	DeploymentID  *int64    `json:"deployment_id"`
	JobType       string    `json:"job_type"`
}

type EventJobStarted struct {
	JobID         int64     `json:"job_id"`
	ServerID      uuid.UUID `json:"server_id"`
	ApplicationID *int64    `json:"application_id"`
	DeploymentID  *int64    `json:"deployment_id"`
	JobType       string    `json:"job_type"`
}

type EventJobFinished struct {
	JobID         int64     `json:"job_id"`
	ServerID      uuid.UUID `json:"server_id"`
	ApplicationID *int64    `json:"application_id"`
	DeploymentID  *int64    `json:"deployment_id"`
	JobType       string    `json:"job_type"`
	Status        JobStatus `json:"status"`
	OutputLog     *string   `json:"output_log"`
}

type EventJobStatusChanged struct {
	JobID  int64     `json:"job_id"`
	Status JobStatus `json:"status"`
}
