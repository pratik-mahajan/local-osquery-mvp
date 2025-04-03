package model

type QueryType string

const (
	QueryTypeOSAndOSQuery QueryType = "os_and_osquery"
	QueryTypeApplications QueryType = "apps"
)

var QueryMap = map[QueryType]string{
	QueryTypeOSAndOSQuery: "SELECT * FROM os_version; SELECT * FROM osquery_info",
	QueryTypeApplications: "SELECT * FROM apps",
}
