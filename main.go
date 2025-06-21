package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/viper"
)

type BookGroup struct {
	Title      string
	Author     string
	Highlights []Highlight
	Expanded   bool
}

type ListItem struct {
	IsBook         bool
	BookIndex      int
	HighlightIndex int
}

type model struct {
	books          []BookGroup
	items          []ListItem
	cursor         int
	selected       map[string]bool // key: "book:highlight" format
	notesDirectory string
	homeDir        string
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

func initialModel() model {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: Could not read config file: %v", err)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Error loading home directory: %v", err)
	}

	notesDirectory := viper.GetString("notes_directory")
	if notesDirectory == "" {
		notesDirectory = "notes"
	}

	highlights, err := parseClippings("My Clippings.txt")
	if err != nil {
		log.Fatalf("Error parsing clippings: %v", err)
	}

	books := groupHighlightsByBook(highlights)
	items := buildItemList(books)

	return model{
		books:          books,
		items:          items,
		selected:       make(map[string]bool),
		notesDirectory: notesDirectory,
		homeDir:        homeDir,
	}
}

func groupHighlightsByBook(highlights []Highlight) []BookGroup {
	bookMap := make(map[string][]Highlight)
	authorMap := make(map[string]string)

	for _, highlight := range highlights {
		bookMap[highlight.Title] = append(bookMap[highlight.Title], highlight)
		authorMap[highlight.Title] = highlight.Author
	}

	var books []BookGroup
	for title, highlights := range bookMap {
		books = append(books, BookGroup{
			Title:      title,
			Author:     authorMap[title],
			Highlights: highlights,
			Expanded:   false,
		})
	}

	return books
}

func buildItemList(books []BookGroup) []ListItem {
	var items []ListItem

	for bookIndex, book := range books {
		items = append(items, ListItem{
			IsBook:    true,
			BookIndex: bookIndex,
		})

		if book.Expanded {
			for highlightIndex := range book.Highlights {
				items = append(items, ListItem{
					IsBook:         false,
					BookIndex:      bookIndex,
					HighlightIndex: highlightIndex,
				})
			}
		}
	}

	return items
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.items)-1 {
				m.cursor++
			}
		case " ":
			currentItem := m.items[m.cursor]
			if currentItem.IsBook {
				// Toggle book expansion
				m.books[currentItem.BookIndex].Expanded = !m.books[currentItem.BookIndex].Expanded
				m.items = buildItemList(m.books)
			} else {
				// Toggle highlight selection
				key := fmt.Sprintf("%d:%d", currentItem.BookIndex, currentItem.HighlightIndex)
				m.selected[key] = !m.selected[key]
			}
		case "a":
			// Select all highlights under current book
			if len(m.items) > 0 {
				currentItem := m.items[m.cursor]
				bookIndex := currentItem.BookIndex
				m.selectAllInBook(bookIndex)
			}
		case "A":
			// Select all highlights globally
			m.selectAllHighlights()
		case "d":
			// Deselect all highlights under current book
			if len(m.items) > 0 {
				currentItem := m.items[m.cursor]
				bookIndex := currentItem.BookIndex
				m.deselectAllInBook(bookIndex)
			}
		case "D":
			// Deselect all highlights globally
			m.deselectAllHighlights()
		case "enter":
			return m, m.exportSelected()
		}
	}
	return m, nil
}

func (m model) View() string {
	s := titleStyle.Render("Kindle Highlights Parser") + "\n\n"
	s += "Navigate: ↑/↓ j/k | Expand/Select: Space | Export: Enter | Quit: q\n"
	s += "Select: a (book) A (all) | Deselect: d (book) D (all)\n\n"

	for i, item := range m.items {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		if item.IsBook {
			// Display book header
			book := m.books[item.BookIndex]
			expandIcon := "▶"
			if book.Expanded {
				expandIcon = "▼"
			}

			// Check how many highlights are selected in this book
			selectedInBook := 0
			for j := range book.Highlights {
				key := fmt.Sprintf("%d:%d", item.BookIndex, j)
				if m.selected[key] {
					selectedInBook++
				}
			}

			bookTitle := book.Title
			if len(bookTitle) > 45 {
				bookTitle = bookTitle[:42] + "..."
			}

			selectionStatus := ""
			if selectedInBook > 0 {
				if selectedInBook == len(book.Highlights) {
					selectionStatus = " [ALL]"
				} else {
					selectionStatus = fmt.Sprintf(" [%d/%d]", selectedInBook, len(book.Highlights))
				}
			}

			line := fmt.Sprintf("%s %s %s (%s) - %d highlights%s",
				cursor, expandIcon, bookTitle, book.Author, len(book.Highlights), selectionStatus)

			if m.cursor == i {
				s += selectedStyle.Render(line) + "\n"
			} else {
				s += bookStyle.Render(line) + "\n"
			}
		} else {
			// Display highlight
			book := m.books[item.BookIndex]
			highlight := book.Highlights[item.HighlightIndex]

			key := fmt.Sprintf("%d:%d", item.BookIndex, item.HighlightIndex)
			checked := " "
			if m.selected[key] {
				checked = "✓"
			}

			text := highlight.Text
			if len(text) > 60 {
				text = text[:57] + "..."
			}

			line := fmt.Sprintf("%s  [%s] %s (Page %s)",
				cursor, checked, text, highlight.Page)

			if m.cursor == i {
				s += selectedStyle.Render(line) + "\n"
			} else {
				s += highlightStyle.Render(line) + "\n"
			}
		}
	}

	selectedCount := len(m.selected)
	s += fmt.Sprintf("\nSelected: %d highlights\n", selectedCount)
	return s
}

func (m model) exportSelected() tea.Cmd {
	return func() tea.Msg {
		bookHighlights := make(map[string][]Highlight)

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

func (m *model) selectAllInBook(bookIndex int) {
	if bookIndex >= len(m.books) {
		return
	}

	book := m.books[bookIndex]
	for i := range book.Highlights {
		key := fmt.Sprintf("%d:%d", bookIndex, i)
		m.selected[key] = true
	}
}

func (m *model) selectAllHighlights() {
	for bookIndex, book := range m.books {
		for i := range book.Highlights {
			key := fmt.Sprintf("%d:%d", bookIndex, i)
			m.selected[key] = true
		}
	}
}

func (m *model) deselectAllInBook(bookIndex int) {
	if bookIndex >= len(m.books) {
		return
	}

	book := m.books[bookIndex]
	for i := range book.Highlights {
		key := fmt.Sprintf("%d:%d", bookIndex, i)
		delete(m.selected, key)
	}
}

func (m *model) deselectAllHighlights() {
	m.selected = make(map[string]bool)
}

func (m model) saveHighlights(title string, highlights []Highlight) error {
	filename := filepath.Join(m.homeDir, m.notesDirectory, title+".md")

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

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v", err)
		os.Exit(1)
	}
}
