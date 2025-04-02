package model

import (
	"time"
)

type OSQueryVersion struct {
	Version     string    `json:"version"`
	ConfigHash  string    `json:"config_hash"`
	ConfigValid int       `json:"config_valid"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type OSQueryVersionSlice []OSQueryVersion

func (s *OSQueryVersionSlice) SetTimestamps(t time.Time) {
	for i := range *s {
		(*s)[i].CreatedAt = t
		(*s)[i].UpdatedAt = t
	}
}
