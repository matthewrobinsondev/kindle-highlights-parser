package tui

import (
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/matthewrobinsdev/kindle-notes-parser/pkg/models"
)

type ExportCompleteMsg struct {
	Results []models.ExportResult
}

func (m *Model) exportSelected() tea.Cmd {
	return func() tea.Msg {
		bookHighlights := make(map[string][]models.Highlight)

		for key, selected := range m.selected {
			if selected {
				// Parse key format "bookIndex:highlightIndex"
				var bookIndex, highlightIndex int
				fmt.Sscanf(key, "%d:%d", &bookIndex, &highlightIndex)

				if bookIndex < len(m.books) && highlightIndex < len(m.books[bookIndex].Highlights) {
					book := m.books[bookIndex]
					highlight := book.Highlights[highlightIndex]
					bookHighlights[book.Title] = append(bookHighlights[book.Title], highlight)
				}
			}
		}

		results, err := m.exporter.ExportHighlights(bookHighlights)
		if err != nil {
			log.Printf("Error exporting highlights: %v", err)
			return ExportCompleteMsg{Results: []models.ExportResult{}}
		}

		return ExportCompleteMsg{Results: results}
	}
}
