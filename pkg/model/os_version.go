package model

import "time"

type OSVersion struct {
	ID           int       `json:"-" db:"id"`
	Name         string    `json:"name" db:"name"`
	Version      string    `json:"version" db:"version"`
	Major        string    `json:"major" db:"major"`
	Minor        string    `json:"minor" db:"minor"`
	Patch        string    `json:"patch" db:"patch"`
	Build        string    `json:"build" db:"build"`
	Platform     string    `json:"platform" db:"platform"`
	PlatformLike string    `json:"platform_like" db:"platform_like"`
	Codename     string    `json:"codename" db:"codename"`
	Arch         string    `json:"arch" db:"arch"`
	Extra        string    `json:"extra" db:"extra"`
	CreatedAt    time.Time `json:"-" db:"created_at"`
	UpdatedAt    time.Time `json:"-" db:"updated_at"`
}
