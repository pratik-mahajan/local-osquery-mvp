package cmdutils

import (
	"encoding/json"
	"fmt"
	"main/pkg/db"
	"main/pkg/model"
	"os/exec"
	"time"
)

type BaseResponse struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ExecuteQuery(queryType model.QueryType) ([]byte, error) {
	query, found := model.QueryMap[queryType]
	if !found {
		return nil, fmt.Errorf("unsupported query type: %s", queryType)
	}

	switch queryType {
	case model.QueryTypeOSAndOSQuery:
		osVersionCmd := exec.Command("osqueryi", "--json", "SELECT * FROM os_version")
		osVersionOutput, err := osVersionCmd.CombinedOutput()
		if err != nil {
			return nil, fmt.Errorf("error running os_version query: %v", err)
		}

		osQueryCmd := exec.Command("osqueryi", "--json", "SELECT * FROM osquery_info")
		osQueryOutput, err := osQueryCmd.CombinedOutput()
		if err != nil {
			return nil, fmt.Errorf("error running osquery_info query: %v", err)
		}

		var osVersions []model.OSVersion
		var osQueryVersions []model.OSQueryVersion

		if err := json.Unmarshal(osVersionOutput, &osVersions); err != nil {
			return nil, fmt.Errorf("error unmarshaling OS version data: %v", err)
		}
		if err := json.Unmarshal(osQueryOutput, &osQueryVersions); err != nil {
			return nil, fmt.Errorf("error unmarshaling OSQuery version data: %v", err)
		}

		if len(osVersions) > 0 && len(osQueryVersions) > 0 {
			if err := db.SaveOSAndOSQueryInfo(osVersions[0], osQueryVersions[0]); err != nil {
				return nil, fmt.Errorf("error saving OS and OSQuery info to database: %v", err)
			}
		}
		combinedOutput := struct {
			OSVersion   []model.OSVersion      `json:"os_version"`
			OSQueryInfo []model.OSQueryVersion `json:"osquery_info"`
		}{
			OSVersion:   osVersions,
			OSQueryInfo: osQueryVersions,
		}

		return json.Marshal(combinedOutput)

	default:
		cmd := exec.Command("osqueryi", "--json", query)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return nil, fmt.Errorf("error running osqueryi: %v", err)
		}

		switch queryType {
		case model.QueryTypeApplications:
			var applications []model.Application
			if err := json.Unmarshal(output, &applications); err != nil {
				return nil, fmt.Errorf("error unmarshaling Applications data: %v", err)
			}
			for _, app := range applications {
				if err := db.SaveApplication(app); err != nil {
					return nil, fmt.Errorf("error saving Application to database: %v", err)
				}
			}
		}

		return output, nil
	}
}

func UnmarshalWithTimestamp(data []byte, v interface{}) error {
	if baseSlice, ok := v.(interface {
		SetTimestamps(time.Time)
	}); ok {
		baseSlice.SetTimestamps(time.Now())
	}

	return json.Unmarshal(data, v)
}

func PrintJSON(v interface{}) (string, error) {
	if data, ok := v.([]byte); ok {
		return string(data), nil
	}
	return "", nil
}
