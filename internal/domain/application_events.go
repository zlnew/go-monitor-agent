package domain

import "github.com/google/uuid"

type EventApplicationCreated struct {
	ApplicationID int64     `json:"application_id"`
	ServerID      uuid.UUID `json:"server_id"`
}

type EventApplicationStatusChanged struct {
	ApplicationID int64             `json:"application_id"`
	Status        ApplicationStatus `json:"status"`
}
