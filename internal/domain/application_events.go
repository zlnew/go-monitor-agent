package domain

type EventApplicationCreated struct {
	Application *Application `json:"application"`
}

type EventApplicationStatusChanged struct {
	ApplicationID int64             `json:"application_id"`
	Status        ApplicationStatus `json:"status"`
	Message       string            `json:"message,omitempty"`
}

type EventApplicationDeploying struct {
	ApplicationID int64  `json:"application_id"`
	Progress      int    `json:"progress"` // 0-100
	Step          string `json:"step"`
}

type EventApplicationDeployed struct {
	ApplicationID int64  `json:"application_id"`
	Success       bool   `json:"success"`
	Message       string `json:"message,omitempty"`
}

type EventApplicationLogs struct {
	ApplicationID int64  `json:"application_id"`
	Logs          string `json:"logs"`
	Timestamp     string `json:"timestamp"`
}

type EventJobStatusChanged struct {
	JobID         int64     `json:"job_id"`
	ApplicationID *int64    `json:"application_id,omitempty"`
	Status        JobStatus `json:"status"`
	OutputLog     *string   `json:"output_log,omitempty"`
}
