package main

import (
	"crispy-doodle/main.go/tui/boba"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {

	if _, err := tea.NewProgram(boba.InitialLogin()).Run(); err != nil {

		fmt.Printf("could not start program: %s\n", err)
		os.Exit(1)
	}
}
