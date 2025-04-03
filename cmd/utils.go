package cmdutils

import (
	"encoding/json"
	"fmt"
	"main/pkg/db"
	"main/pkg/model"
	"os/exec"
	"time"
)

type OutputFormat int

const (
	FormatJSON OutputFormat = iota
	FormatText
)

func ExecuteQuery(queryType model.QueryType, format OutputFormat) ([]byte, error) {
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

			timestamp := time.Now().Format("2006-01-02 15:04:05")

			if format == FormatJSON {
				combinedData := struct {
					OSVersion   []model.OSVersion      `json:"os_version"`
					OSQueryInfo []model.OSQueryVersion `json:"osquery_info"`
				}{
					OSVersion:   osVersions,
					OSQueryInfo: osQueryVersions,
				}

				jsonData, err := json.Marshal(combinedData)
				if err != nil {
					return nil, fmt.Errorf("error marshaling combined data: %v", err)
				}
				return jsonData, nil
			}

			statusMsg := fmt.Sprintf("🕒 [%s] ✅ OS (%s %s) and OSQuery (v%s) info stored in database\n",
				timestamp,
				osVersions[0].Name,
				osVersions[0].Version,
				osQueryVersions[0].Version)
			return []byte(statusMsg), nil
		}
		if format == FormatJSON {
			return []byte("[]"), nil
		}
		return []byte("No data found"), nil

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

			count := 0
			for _, app := range applications {
				if err := db.SaveApplication(app); err != nil {
					return nil, fmt.Errorf("error saving Application to database: %v", err)
				}
				count++
			}

			timestamp := time.Now().Format("2006-01-02 15:04:05")

			if format == FormatJSON {
				return output, nil
			}

			statusMsg := fmt.Sprintf("🕒 [%s] ✅ %d applications stored in database\n", timestamp, count)
			return []byte(statusMsg), nil
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
