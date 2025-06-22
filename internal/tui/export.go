package tui

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

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

		var results []models.ExportResult
		for title, highlights := range bookHighlights {
			result, err := m.saveHighlightsWithDuplicateCheck(title, highlights)
			if err != nil {
				log.Printf("Error saving highlights for %s: %v", title, err)
				continue
			}
			results = append(results, result)
		}

		return ExportCompleteMsg{Results: results}
	}
}

func (m *Model) saveHighlightsWithDuplicateCheck(title string, highlights []models.Highlight) (models.ExportResult, error) {
	filename := filepath.Join(m.config.HomeDir, m.config.NotesDirectory, title+".md")

	if err := os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
		return models.ExportResult{}, fmt.Errorf("error creating directory: %w", err)
	}

	existingHighlights := make(map[string]bool)
	var existingContent strings.Builder

	if _, err := os.Stat(filename); err == nil {
		content, err := os.ReadFile(filename)
		if err != nil {
			return models.ExportResult{}, fmt.Errorf("error reading existing file: %w", err)
		}

		existingContent.Write(content)
		existingHighlights = extractExistingHighlights(string(content))
	} else {
		existingContent.WriteString(fmt.Sprintf("# %s\n\n", title))
	}

	var newHighlights []models.Highlight
	skippedCount := 0

	for _, highlight := range highlights {
		highlightKey := createHighlightKey(highlight)
		if !existingHighlights[highlightKey] {
			newHighlights = append(newHighlights, highlight)
		} else {
			skippedCount++
		}
	}

	for _, highlight := range newHighlights {
		existingContent.WriteString(fmt.Sprintf("- %s (Page: %s)\n", highlight.Text, highlight.Page))
	}

	if len(newHighlights) > 0 {
		if err := os.WriteFile(filename, []byte(existingContent.String()), 0644); err != nil {
			return models.ExportResult{}, fmt.Errorf("error writing file: %w", err)
		}
	}

	return models.ExportResult{
		BookTitle:    title,
		NewCount:     len(newHighlights),
		SkippedCount: skippedCount,
		TotalCount:   len(highlights),
	}, nil
}

func extractExistingHighlights(content string) map[string]bool {
	highlights := make(map[string]bool)

	// Regex to match highlight lines: "- text (Page: number)"
	re := regexp.MustCompile(`^- (.+) \(Page: (\d+)\)$`)

	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if matches := re.FindStringSubmatch(line); len(matches) == 3 {
			text := matches[1]
			page := matches[2]
			key := fmt.Sprintf("%s|%s", text, page)
			highlights[key] = true
		}
	}

	return highlights
}

func createHighlightKey(highlight models.Highlight) string {
	return fmt.Sprintf("%s|%s", highlight.Text, highlight.Page)
}
