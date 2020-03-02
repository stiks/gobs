package models

import "errors"

var (
	// ErrMissCache ...
	ErrMissCache = errors.New("cache not found")
	// ErrCache ...
	ErrCache = errors.New("cache error")
)
