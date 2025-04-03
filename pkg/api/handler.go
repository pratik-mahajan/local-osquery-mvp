package api

import (
	"encoding/json"
	"fmt"
	cmdutils "main/cmd"
	"main/pkg/model"
	"net/http"
	"time"
)

type LatestDataResponse struct {
	OSVersion    *model.OSVersion      `json:"os_version"`
	OSQueryInfo  *model.OSQueryVersion `json:"osquery_info"`
	Applications []model.Application   `json:"applications"`
	Errors       []string              `json:"errors,omitempty"`
}

type combinedOSData struct {
	OSVersion   []model.OSVersion      `json:"os_version"`
	OSQueryInfo []model.OSQueryVersion `json:"osquery_info"`
}

func logRequest(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		handler(w, r)
		duration := time.Since(start)
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		fmt.Printf("🌐 [%s] %s %s (took %v)\n", timestamp, r.Method, r.URL.Path, duration)
	}
}

func RunServer(port string) error {
	http.HandleFunc("/latest_data", logRequest(handleLatestData))
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
		model.QueryTypeOSAndOSQuery,
		model.QueryTypeApplications,
	}

	for _, queryType := range queryTypes {
		output, err := cmdutils.ExecuteQuery(queryType, cmdutils.FormatJSON)
		if err != nil {
			errors = append(errors, err.Error())
			continue
		}

		switch queryType {
		case model.QueryTypeOSAndOSQuery:
			var combinedData combinedOSData
			if err := json.Unmarshal(output, &combinedData); err != nil {
				errors = append(errors, "Error parsing OS and OSQuery data: "+err.Error())
				continue
			}

			if len(combinedData.OSVersion) > 0 {
				response.OSVersion = &combinedData.OSVersion[0]
			}
			if len(combinedData.OSQueryInfo) > 0 {
				response.OSQueryInfo = &combinedData.OSQueryInfo[0]
			}

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
