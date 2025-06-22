package tui

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/matthewrobinsdev/kindle-notes-parser/internal/config"
	"github.com/matthewrobinsdev/kindle-notes-parser/internal/exporter"
	"github.com/matthewrobinsdev/kindle-notes-parser/internal/parser"
	"github.com/matthewrobinsdev/kindle-notes-parser/pkg/models"
)

type Model struct {
	books    []models.BookGroup
	items    []models.ListItem
	cursor   int
	selected map[string]bool // key: "book:highlight" format
	config   *config.Config
	exporter *exporter.Service
}

var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1)

	bookStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4")).
			Bold(true)

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#EE6FF8"))

	normalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA"))

	highlightStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#CCCCCC")).
			PaddingLeft(2)
)

func NewModel(clippingsFile string) *Model {
	cfg := config.Load()

	highlights, err := parser.ParseClippings(clippingsFile)
	if err != nil {
		log.Fatalf("Error parsing clippings: %v", err)
	}

	books := parser.GroupHighlightsByBook(highlights)
	items := buildItemList(books)

	return &Model{
		books:    books,
		items:    items,
		selected: make(map[string]bool),
		config:   cfg,
		exporter: exporter.New(cfg),
	}
}

func buildItemList(books []models.BookGroup) []models.ListItem {
	var items []models.ListItem

	for bookIndex, book := range books {
		items = append(items, models.ListItem{
			IsBook:    true,
			BookIndex: bookIndex,
		})

		if book.Expanded {
			for highlightIndex := range book.Highlights {
				items = append(items, models.ListItem{
					IsBook:         false,
					BookIndex:      bookIndex,
					HighlightIndex: highlightIndex,
				})
			}
		}
	}

	return items
}

func (m *Model) Init() tea.Cmd {
	return nil
}
