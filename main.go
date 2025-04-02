package main

import (
	"fmt"
	cmdutils "main/cmd"
	"main/pkg/db"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	if err := db.InitDB(); err != nil {
		fmt.Printf("Error initializing database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	p := tea.NewProgram(cmdutils.NewMenuModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
