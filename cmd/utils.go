package cmdutils

import (
	"encoding/json"
	"fmt"
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

	return output, nil
}

func UnmarshalWithTimestamp(data []byte, v interface{}) error {
	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("error parsing JSON output: %v", err)
	}

	if baseSlice, ok := v.(interface {
		SetTimestamps(time.Time)
	}); ok {
		baseSlice.SetTimestamps(time.Now())
	}

	return nil
}

func PrintJSON(v interface{}) (string, error) {
	prettyJSON, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", fmt.Errorf("error formatting JSON: %v", err)
	}
	return string(prettyJSON), nil
}
