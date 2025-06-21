package tui

import (
	"fmt"
)

func (m *Model) View() string {
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
