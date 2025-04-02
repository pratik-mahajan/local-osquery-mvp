package model

import (
	"time"
)

type Application struct {
	Name                 string    `json:"name"`
	Path                 string    `json:"path"`
	BundleExecutable     string    `json:"bundle_executable"`
	BundleIdentifier     string    `json:"bundle_identifier"`
	BundleName           string    `json:"bundle_name"`
	BundleShortVersion   string    `json:"bundle_short_version"`
	BundleVersion        string    `json:"bundle_version"`
	BundlePackageType    string    `json:"bundle_package_type"`
	Environment          string    `json:"environment"`
	Element              string    `json:"element"`
	Compiler             string    `json:"compiler"`
	DevelopmentRegion    string    `json:"development_region"`
	DisplayName          string    `json:"display_name"`
	InfoString           string    `json:"info_string"`
	MinimumSystemVersion string    `json:"minimum_system_version"`
	Category             string    `json:"category"`
	ApplescriptEnabled   string    `json:"applescript_enabled"`
	Copyright            string    `json:"copyright"`
	LastOpenedTime       string `json:"last_opened_time"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

type ApplicationSlice []Application

func (s *ApplicationSlice) SetTimestamps(t time.Time) {
	for i := range *s {
		(*s)[i].CreatedAt = t
		(*s)[i].UpdatedAt = t
	}
}
