package exporter

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/matthewrobinsdev/kindle-notes-parser/internal/config"
	"github.com/matthewrobinsdev/kindle-notes-parser/pkg/models"
)

type MockFileSystem struct {
	files map[string][]byte
	dirs  map[string]bool
}

func NewMockFileSystem() *MockFileSystem {
	return &MockFileSystem{
		files: make(map[string][]byte),
		dirs:  make(map[string]bool),
	}
}

func (fs *MockFileSystem) Stat(name string) (os.FileInfo, error) {
	if _, exists := fs.files[name]; exists {
		return nil, nil // File exists
	}
	return nil, os.ErrNotExist
}

func (fs *MockFileSystem) ReadFile(filename string) ([]byte, error) {
	if data, exists := fs.files[filename]; exists {
		return data, nil
	}
	return nil, os.ErrNotExist
}

func (fs *MockFileSystem) WriteFile(filename string, data []byte, perm os.FileMode) error {
	fs.files[filename] = data
	return nil
}

func (fs *MockFileSystem) MkdirAll(path string, perm os.FileMode) error {
	fs.dirs[path] = true
	return nil
}

func TestExtractExistingHighlights(t *testing.T) {
	tests := []struct {
		name           string
		content        string
		expectedCount  int
		expectedKeys   []string
		unexpectedKeys []string
	}{
		{
			name: "multiple valid highlights",
			content: `# Test Book

- First highlight text (Page: 123)
- Second highlight text (Page: 456)
- Third highlight with special chars! (Page: 789)
`,
			expectedCount: 3,
			expectedKeys: []string{
				"First highlight text|123",
				"Second highlight text|456",
				"Third highlight with special chars!|789",
			},
			unexpectedKeys: []string{"Non-existent highlight|999"},
		},
		{
			name: "empty content",
			content: `# Empty Book

`,
			expectedCount:  0,
			expectedKeys:   []string{},
			unexpectedKeys: []string{"Any highlight|123"},
		},
		{
			name: "mixed content with non-matching lines",
			content: `# Test Book

Some random text
- Not a highlight line
- Almost a highlight (Page: )
- Valid highlight (Page: 123)
Another random line
`,
			expectedCount:  1,
			expectedKeys:   []string{"Valid highlight|123"},
			unexpectedKeys: []string{"Not a highlight line|123", "Almost a highlight|"},
		},
		{
			name: "highlights with complex text",
			content: `# Complex Book

- Highlight with "quotes" and symbols! @#$% (Page: 42)
- Multi-word highlight with numbers 123 and punctuation... (Page: 999)
`,
			expectedCount: 2,
			expectedKeys: []string{
				"Highlight with \"quotes\" and symbols! @#$%|42",
				"Multi-word highlight with numbers 123 and punctuation...|999",
			},
			unexpectedKeys: []string{},
		},
		{
			name: "no highlights section",
			content: `# Book Title

This is just regular text.
No highlights here.
`,
			expectedCount:  0,
			expectedKeys:   []string{},
			unexpectedKeys: []string{"This is just regular text.|123"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{HomeDir: "/test", NotesDirectory: "notes"}
			service := New(cfg)
			existing := service.extractExistingHighlights(tt.content)

			assert.Equal(t, tt.expectedCount, len(existing), "Should extract correct number of highlights")

			for _, key := range tt.expectedKeys {
				assert.True(t, existing[key], "Should find expected key: %s", key)
			}

			for _, key := range tt.unexpectedKeys {
				assert.False(t, existing[key], "Should not find unexpected key: %s", key)
			}
		})
	}
}

func TestCreateHighlightKey(t *testing.T) {
	tests := []struct {
		name      string
		highlight models.Highlight
		expected  string
	}{
		{
			name: "simple highlight",
			highlight: models.Highlight{
				Text: "Test highlight text",
				Page: "123",
			},
			expected: "Test highlight text|123",
		},
		{
			name: "highlight with special characters",
			highlight: models.Highlight{
				Text: "Text with \"quotes\" and symbols! @#$%",
				Page: "42",
			},
			expected: "Text with \"quotes\" and symbols! @#$%|42",
		},
		{
			name: "highlight with numbers and punctuation",
			highlight: models.Highlight{
				Text: "Numbers 123 and punctuation...",
				Page: "999",
			},
			expected: "Numbers 123 and punctuation...|999",
		},
		{
			name: "empty text",
			highlight: models.Highlight{
				Text: "",
				Page: "1",
			},
			expected: "|1",
		},
		{
			name: "single character",
			highlight: models.Highlight{
				Text: "A",
				Page: "5",
			},
			expected: "A|5",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{HomeDir: "/test", NotesDirectory: "notes"}
			service := New(cfg)
			result := service.createHighlightKey(tt.highlight)
			assert.Equal(t, tt.expected, result, "Should create correct highlight key")
		})
	}
}

func TestSanitizeFilename(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "normal filename",
			input:    "Normal Book Title",
			expected: "Normal Book Title",
		},
		{
			name:     "filename with problematic characters",
			input:    "Book/Title\\With:*?\"<>|Chars",
			expected: "Book-Title-With--Chars",
		},
		{
			name:     "filename with whitespace",
			input:    "  Book Title  ",
			expected: "Book Title",
		},
		{
			name:     "very long filename",
			input:    strings.Repeat("A", 250),
			expected: strings.Repeat("A", 200),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{HomeDir: "/test", NotesDirectory: "notes"}
			service := New(cfg)
			result := service.sanitizeFilename(tt.input)
			assert.Equal(t, tt.expected, result, "Should sanitize filename correctly")
		})
	}
}

func TestExportHighlightsWithMockFS(t *testing.T) {
	cfg := &config.Config{
		HomeDir:        "/home/user",
		NotesDirectory: "notes",
	}

	mockFS := NewMockFileSystem()
	service := NewWithFileSystem(cfg, mockFS)

	highlights := map[string][]models.Highlight{
		"Test Book": {
			{Text: "First highlight", Page: "1"},
			{Text: "Second highlight", Page: "2"},
		},
	}

	results, err := service.ExportHighlights(highlights)

	require.NoError(t, err, "Should export without error")
	require.Len(t, results, 1, "Should have one result")

	result := results[0]
	assert.Equal(t, "Test Book", result.BookTitle)
	assert.Equal(t, 2, result.NewCount)
	assert.Equal(t, 0, result.SkippedCount)
	assert.Equal(t, 2, result.TotalCount)

	expectedPath := "/home/user/notes/Test Book.md"
	data, exists := mockFS.files[expectedPath]
	require.True(t, exists, "File should be created")

	content := string(data)
	assert.Contains(t, content, "# Test Book")
	assert.Contains(t, content, "- First highlight (Page: 1)")
	assert.Contains(t, content, "- Second highlight (Page: 2)")
}
