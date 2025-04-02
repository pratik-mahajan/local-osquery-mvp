package main

import (
	"flag"
	"fmt"
	cmdutils "main/cmd"
	"main/pkg/api"
	"main/pkg/db"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	serverMode := flag.Bool("server", false, "Run in server mode")
	port := flag.String("port", "8080", "Port to run the server on")
	flag.Parse()

	if err := db.InitDB(); err != nil {
		fmt.Printf("Error initializing database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	if *serverMode {
		fmt.Printf("\n🚀 Server is running at http://localhost:%s\n", *port)
		fmt.Printf("\n💡 Try it out: ")
		fmt.Printf("   curl http://localhost:%s/latest_data\n\n", *port)
		if err := api.RunServer(*port); err != nil {
			fmt.Printf("Error running server: %v\n", err)
			os.Exit(1)
		}
	} else {
		p := tea.NewProgram(cmdutils.NewMenuModel())
		if _, err := p.Run(); err != nil {
			fmt.Printf("Error running program: %v\n", err)
			os.Exit(1)
		}
	}
}
