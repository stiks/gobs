package models

import "time"

type Stats struct {
	Uptime       time.Time      `json:"uptime"`
	RequestCount uint64         `json:"requestCount"`
	Statuses     map[string]int `json:"statuses"`
}
