package main

import (
	"encoding/json"
	"fmt"
	"log"
	"main/pkg/model"
	"os/exec"
	"time"
)

func main() {
	cmd := exec.Command("osqueryi", "--json", "SELECT * FROM os_version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Error running osqueryi: %v", err)
	}

	var versions []model.OSVersion
	if err := json.Unmarshal(output, &versions); err != nil {
		log.Fatalf("Error parsing JSON output: %v", err)
	}

	now := time.Now()
	for i := range versions {
		versions[i].CreatedAt = now
		versions[i].UpdatedAt = now
	}

	prettyJSON, err := json.MarshalIndent(versions, "", "  ")
	if err != nil {
		log.Fatalf("Error formatting JSON: %v", err)
	}
	fmt.Printf("OS Version Data:\n%s\n", string(prettyJSON))
}
