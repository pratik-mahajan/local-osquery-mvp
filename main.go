package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
)

type OSVersion struct {
	ID           int       `json:"-" db:"id"`
	Name         string    `json:"name" db:"name"`
	Version      string    `json:"version" db:"version"`
}

func main() {
	cmd := exec.Command("osqueryi", "--json", "SELECT * FROM os_version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Error running osqueryi: %v", err)
	}

	var versions []OSVersion
	if err := json.Unmarshal(output, &versions); err != nil {
		log.Fatalf("Error parsing JSON output: %v", err)
	}

	prettyJSON, err := json.MarshalIndent(versions, "", "  ")
	if err != nil {
		log.Fatalf("Error formatting JSON: %v", err)
	}
	fmt.Printf("OS Version Data:\n%s\n", string(prettyJSON))
}
