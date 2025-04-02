package model

type QueryType string

const (
	QueryTypeOSVersion      QueryType = "os_version"
	QueryTypeOSQueryVersion QueryType = "osquery_info"
	QueryTypeApplications   QueryType = "apps"
)

var QueryMap = map[QueryType]string{
	QueryTypeOSVersion:      "SELECT * FROM os_version",
	QueryTypeOSQueryVersion: "SELECT * FROM osquery_info",
	QueryTypeApplications:   "SELECT name, path, bundle_short_version FROM apps",
}
