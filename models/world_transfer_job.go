package models

import (
	"time"

	"github.com/google/uuid"
)

type WorldTransferJobStatus string

const (
	WorldTransferJobStatusCreated   WorldTransferJobStatus = "created"
	WorldTransferJobStatusCompleted WorldTransferJobStatus = "completed"
)

type WorldsTransferJob struct {
	ID                uuid.UUID              `json:"id"`
	TargetEnvironment string                 `json:"target_environment"`
	UserID            uuid.UUID              `json:"user_id"`
	Status            WorldTransferJobStatus `json:"status"`
	CreatedAt         time.Time              `json:"created_at"`
	UpdatedAt         time.Time              `json:"updated_at"`
}

type WorldTransferJob struct {
	JobId        uuid.UUID              `json:"job_id"`
	WorldID      uuid.UUID              `json:"world_id"`
	WorldVersion int                    `json:"world_version"`
	Status       WorldTransferJobStatus `json:"status"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
}

type WorldTransferJobStatusDTO struct {
	JobId           uuid.UUID                            `json:"job_id"`
	Status          WorldTransferJobStatus               `json:"status"`
	StatusByWorldID map[uuid.UUID]WorldTransferJobStatus `json:"status_by_world_id"`
}
