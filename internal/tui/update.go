package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// Update handles messages and updates the model
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m *Model) selectAllInBook(bookIndex int) {
	if bookIndex >= len(m.books) {
		return
	}

	book := m.books[bookIndex]
	for i := range book.Highlights {
		key := fmt.Sprintf("%d:%d", bookIndex, i)
		m.selected[key] = true
	}
}

func (m *Model) selectAllHighlights() {
	for bookIndex, book := range m.books {
		for i := range book.Highlights {
			key := fmt.Sprintf("%d:%d", bookIndex, i)
			m.selected[key] = true
		}
	}
}

func (m *Model) deselectAllInBook(bookIndex int) {
	if bookIndex >= len(m.books) {
		return
	}

	book := m.books[bookIndex]
	for i := range book.Highlights {
		key := fmt.Sprintf("%d:%d", bookIndex, i)
		delete(m.selected, key)
	}
}

func (m *Model) deselectAllHighlights() {
	m.selected = make(map[string]bool)
}
