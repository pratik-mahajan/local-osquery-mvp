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

	cmd := exec.Command("osqueryi", "--json", query)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("error running osqueryi: %v", err)
	}

	switch queryType {
	case model.QueryTypeOSVersion:
		var osVersions []model.OSVersion
		if err := json.Unmarshal(output, &osVersions); err != nil {
			return nil, fmt.Errorf("error unmarshaling OS version data: %v", err)
		}
		if len(osVersions) > 0 {
			if err := db.SaveOSVersion(osVersions[0]); err != nil {
				return nil, fmt.Errorf("error saving OS version to database: %v", err)
			}
		}
	case model.QueryTypeOSQueryVersion:
		var osQueryVersions []model.OSQueryVersion
		if err := json.Unmarshal(output, &osQueryVersions); err != nil {
			return nil, fmt.Errorf("error unmarshaling OSQuery version data: %v", err)
		}
		if len(osQueryVersions) > 0 {
			if err := db.SaveOSQueryVersion(osQueryVersions[0]); err != nil {
				return nil, fmt.Errorf("error saving OSQuery version to database: %v", err)
			}
		}
	}

	return output, nil
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
