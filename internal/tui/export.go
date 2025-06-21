package tui

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/matthewrobinsdev/kindle-notes-parser/pkg/models"
)

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

		for title, highlights := range bookHighlights {
			if err := m.saveHighlights(title, highlights); err != nil {
				log.Printf("Error saving highlights for %s: %v", title, err)
			}
		}

		return tea.Quit()
	}
}

func (m *Model) saveHighlights(title string, highlights []models.Highlight) error {
	filename := filepath.Join(m.config.HomeDir, m.config.NotesDirectory, title+".md")

	if err := os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
		return err
	}

	var content strings.Builder
	content.WriteString(fmt.Sprintf("# %s\n\n", title))

	for _, highlight := range highlights {
		content.WriteString(fmt.Sprintf("- %s (Page: %s)\n", highlight.Text, highlight.Page))
	}

	return os.WriteFile(filename, []byte(content.String()), 0644)
}
