package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/matthewrobinsdev/kindle-notes-parser/internal/tui"
)

func main() {
	clippingsFile := "My Clippings.txt"

	if _, err := os.Stat(clippingsFile); os.IsNotExist(err) {
		fmt.Printf("Error: %s not found. Please place your Kindle clippings file in the current directory.\n", clippingsFile)
		os.Exit(1)
	}

	model := tui.NewModel(clippingsFile)
	p := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v", err)
		os.Exit(1)
	}
}
