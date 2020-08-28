package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	// ErrClientNotFound ...
	ErrClientNotFound = errors.New("client not found")
	// ErrClientNameTaken ...
	ErrClientNameTaken = errors.New("client name taken")
)

// ClientQueryParams ...
type ClientQueryParams struct {
	Page    int    `query:"current"`
	PerPage int    `query:"pageSize"`
	Role    string `query:"role"`
	Status  *int   `query:"status"`
	Query   string `query:"query"`
}

// Client model
type Client struct {
	ID        uuid.UUID `json:"id"         sql:"type:uuid,pk"`
	Name      string    `json:"firstName"  sql:"type:varchar(255)"`
	Email     string    `json:"email"      sql:",unique,index"`
	Status    int       `json:"status"`
	OwnerID   uuid.UUID `json:"ownerId"    sql:",type:uuid"`
	Owner     *Client   `json:"owner"`
	CreatedAt time.Time `json:"createdAt"  sql:"default:now()"`
	UpdatedAt time.Time `json:"updatedAt"  sql:"default:now()"`
}

// CreateClient model
type CreateClient struct {
	Name   string `json:"firstName"  sql:"type:varchar(255)"`
	Email  string `json:"email"      sql:",unique,index"`
	Status int    `json:"status"`
}
