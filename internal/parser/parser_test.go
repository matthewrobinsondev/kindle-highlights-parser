package parser

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/matthewrobinsdev/kindle-notes-parser/pkg/models"
)

const CLIPPINGS_FILE_PATH = "../../testData/Test Clippings.txt"
const STRANGE_CLIPPINGS_FILE_PATH = "../../testData/Strange Title Clippings.txt"
const FORMATTED_MARKDOWN_FILE_PATH = "../../testData/SandwormFormatted.md"

func TestParseClippings(t *testing.T) {
	tests := []struct {
		name                   string
		filePath               string
		expectedBookCounts     map[string]int
		expectedFirstHighlight *models.Highlight
		shouldError            bool
	}{
		{
			name:     "standard clippings file",
			filePath: CLIPPINGS_FILE_PATH,
			expectedBookCounts: map[string]int{
				"Sandworm":                    2,
				"Modern Software Engineering": 1,
			},
			expectedFirstHighlight: &models.Highlight{
				Title:  "Sandworm",
				Author: "Greenberg, Andy",
				Text:   "Put more simply, a complex system like a digitized civilization is subject to cascading failures, where one thing depends on another, which depends on another thing.",
			},
			shouldError: false,
		},
		{
			name:     "strange book name file",
			filePath: STRANGE_CLIPPINGS_FILE_PATH,
			expectedBookCounts: map[string]int{
				"Sandworm": 2,
			},
			expectedFirstHighlight: nil, // We don't care about specific content for this test
			shouldError:            false,
		},
		{
			name:                   "non-existent file",
			filePath:               "non-existent-file.txt",
			expectedBookCounts:     nil,
			expectedFirstHighlight: nil,
			shouldError:            true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.filePath != "non-existent-file.txt" {
				require.FileExists(t, tt.filePath, "Test file should exist")
			}

			highlights, err := ParseClippings(tt.filePath)

			if tt.shouldError {
				assert.Error(t, err, "Should return error for invalid file")
				return
			}

			require.NoError(t, err, "Should not return error for valid file")

			bookCounts := make(map[string]int)
			for _, highlight := range highlights {
				bookCounts[highlight.Title]++
			}

			for expectedBook, expectedCount := range tt.expectedBookCounts {
				assert.Equal(t, expectedCount, bookCounts[expectedBook],
					"Should have correct number of highlights for %s", expectedBook)
			}

			if tt.expectedFirstHighlight != nil {
				var foundHighlight *models.Highlight
				for _, highlight := range highlights {
					if highlight.Title == tt.expectedFirstHighlight.Title {
						foundHighlight = &highlight
						break
					}
				}

				require.NotNil(t, foundHighlight, "Should find expected highlight")
				assert.Equal(t, tt.expectedFirstHighlight.Text, foundHighlight.Text, "Highlight text should match")
				assert.Equal(t, tt.expectedFirstHighlight.Author, foundHighlight.Author, "Author should match")
			}
		})
	}
}

func TestGroupHighlightsByBook(t *testing.T) {
	tests := []struct {
		name               string
		highlights         []models.Highlight
		expectedBookCount  int
		expectedBookTitles []string
	}{
		{
			name: "multiple books with multiple highlights",
			highlights: []models.Highlight{
				{Title: "Book A", Author: "Author 1", Text: "Highlight 1", Page: "1"},
				{Title: "Book A", Author: "Author 1", Text: "Highlight 2", Page: "2"},
				{Title: "Book B", Author: "Author 2", Text: "Highlight 3", Page: "3"},
			},
			expectedBookCount:  2,
			expectedBookTitles: []string{"Book A", "Book B"},
		},
		{
			name: "single book with single highlight",
			highlights: []models.Highlight{
				{Title: "Solo Book", Author: "Solo Author", Text: "Solo Highlight", Page: "1"},
			},
			expectedBookCount:  1,
			expectedBookTitles: []string{"Solo Book"},
		},
		{
			name:               "empty highlights",
			highlights:         []models.Highlight{},
			expectedBookCount:  0,
			expectedBookTitles: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			books := GroupHighlightsByBook(tt.highlights)

			assert.Equal(t, tt.expectedBookCount, len(books), "Should have correct number of books")

			bookTitles := make([]string, len(books))
			for i, book := range books {
				bookTitles[i] = book.Title
			}

			for _, expectedTitle := range tt.expectedBookTitles {
				assert.Contains(t, bookTitles, expectedTitle, "Should contain expected book title")
			}

			if len(books) > 0 {
				for _, book := range books {
					assert.NotEmpty(t, book.Title, "Book should have a title")
					assert.NotEmpty(t, book.Author, "Book should have an author")
					assert.NotEmpty(t, book.Highlights, "Book should have highlights")
					assert.False(t, book.Expanded, "Book should start collapsed")

					for _, highlight := range book.Highlights {
						assert.Equal(t, book.Title, highlight.Title, "All highlights should have same title as book")
					}
				}
			}
		})
	}
}

func TestFormatMarkdown(t *testing.T) {
	require.FileExists(t, CLIPPINGS_FILE_PATH, "Test clippings file should exist")
	require.FileExists(t, FORMATTED_MARKDOWN_FILE_PATH, "Expected markdown file should exist")

	highlights, err := ParseClippings(CLIPPINGS_FILE_PATH)
	require.NoError(t, err, "Should parse clippings without error")

	var sandwormHighlights []models.Highlight
	for _, highlight := range highlights {
		if highlight.Title == "Sandworm" {
			sandwormHighlights = append(sandwormHighlights, highlight)
		}
	}

	require.NotEmpty(t, sandwormHighlights, "Should find Sandworm highlights")

	result := formatMarkdownForFile(sandwormHighlights)

	expectedContent, err := os.ReadFile(FORMATTED_MARKDOWN_FILE_PATH)
	require.NoError(t, err, "Should read expected markdown file")

	assert.Equal(t, string(expectedContent), result, "Formatted markdown should match expected output")
}

func formatMarkdownForFile(highlights []models.Highlight) string {
	var markdown strings.Builder

	for _, highlight := range highlights {
		markdown.WriteString(fmt.Sprintf("- %s (Page: %s)\n", highlight.Text, highlight.Page))
	}

	return markdown.String()
}
