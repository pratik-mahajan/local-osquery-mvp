package model

import (
	"time"
)

type OSQueryVersion struct {
	PID           int       `json:"pid"`
	UUID          string    `json:"uuid"`
	InstanceID    string    `json:"instance_id"`
	Version       string    `json:"version"`
	ConfigHash    string    `json:"config_hash"`
	ConfigValid   string    `json:"config_valid"`
	Extensions    string    `json:"extensions"`
	BuildPlatform string    `json:"build_platform"`
	BuildDistro   string    `json:"build_distro"`
	StartTime     int       `json:"start_time"`
	Watcher       int       `json:"watcher"`
	PlatformMask  int       `json:"platform_mask"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type OSQueryVersionSlice []OSQueryVersion

func (s *OSQueryVersionSlice) SetTimestamps(t time.Time) {
	for i := range *s {
		(*s)[i].CreatedAt = t
		(*s)[i].UpdatedAt = t
	}
}
