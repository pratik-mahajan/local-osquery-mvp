package model

import (
	"time"
)

type Application struct {
	Name      string    `json:"name"`
	Path      string    `json:"path"`
	Version   string    `json:"bundle_short_version"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ApplicationSlice []Application

func (s *ApplicationSlice) SetTimestamps(t time.Time) {
	for i := range *s {
		(*s)[i].CreatedAt = t
		(*s)[i].UpdatedAt = t
	}
}
