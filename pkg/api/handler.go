package api

import (
	"encoding/json"
	cmdutils "main/cmd"
	"main/pkg/model"
	"net/http"
)

type LatestDataResponse struct {
	OSVersion    []model.OSVersion      `json:"os_version"`
	OSQueryInfo  []model.OSQueryVersion `json:"osquery_info"`
	Applications []model.Application    `json:"applications"`
	Errors       []string               `json:"errors,omitempty"`
}

func RunServer(port string) error {
	http.HandleFunc("/latest_data", handleLatestData)
	return http.ListenAndServe(":"+port, nil)
}

func handleLatestData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response := LatestDataResponse{}
	var errors []string

	queryTypes := []model.QueryType{
		model.QueryTypeOSVersion,
		model.QueryTypeOSQueryVersion,
		model.QueryTypeApplications,
	}

	for _, queryType := range queryTypes {
		output, err := cmdutils.ExecuteQuery(queryType)
		if err != nil {
			errors = append(errors, err.Error())
			continue
		}

		switch queryType {
		case model.QueryTypeOSVersion:
			var osVersions []model.OSVersion
			if err := json.Unmarshal(output, &osVersions); err != nil {
				errors = append(errors, "Error parsing OS version data: "+err.Error())
				continue
			}
			response.OSVersion = osVersions

		case model.QueryTypeOSQueryVersion:
			var osQueryVersions []model.OSQueryVersion
			if err := json.Unmarshal(output, &osQueryVersions); err != nil {
				errors = append(errors, "Error parsing OSQuery info data: "+err.Error())
				continue
			}
			response.OSQueryInfo = osQueryVersions

		case model.QueryTypeApplications:
			var applications []model.Application
			if err := json.Unmarshal(output, &applications); err != nil {
				errors = append(errors, "Error parsing Applications data: "+err.Error())
				continue
			}
			response.Applications = applications
		}
	}

	if len(errors) > 0 {
		response.Errors = errors
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
