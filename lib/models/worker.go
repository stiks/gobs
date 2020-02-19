package models

import "github.com/google/uuid"

// WorkerRequest ...
type WorkerRequest struct {
	ID uuid.UUID `json:"id"`
}
