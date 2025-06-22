package tui

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/matthewrobinsdev/kindle-notes-parser/pkg/models"
)

func TestExtractExistingHighlights(t *testing.T) {
	assert := assert.New(t)
	
	content := `# Test Book

- First highlight text (Page: 123)
- Second highlight text (Page: 456)
- Third highlight with special chars! (Page: 789)
`
	
	existing := extractExistingHighlights(content)
	
	assert.Equal(3, len(existing), "Should extract 3 highlights")
	assert.True(existing["First highlight text|123"], "Should find first highlight")
	assert.True(existing["Second highlight text|456"], "Should find second highlight")
	assert.True(existing["Third highlight with special chars!|789"], "Should find third highlight")
	assert.False(existing["Non-existent highlight|999"], "Should not find non-existent highlight")
}

func TestCreateHighlightKey(t *testing.T) {
	assert := assert.New(t)
	
	highlight := models.Highlight{
		Text: "Test highlight text",
		Page: "123",
	}
	
	key := createHighlightKey(highlight)
	expected := "Test highlight text|123"
	
	assert.Equal(expected, key, "Should create correct highlight key")
}

func TestExtractExistingHighlightsEmptyContent(t *testing.T) {
	assert := assert.New(t)
	
	content := `# Empty Book

`
	
	existing := extractExistingHighlights(content)
	
	assert.Equal(0, len(existing), "Should extract 0 highlights from empty content")
}

func TestExtractExistingHighlightsWithNonMatchingLines(t *testing.T) {
	assert := assert.New(t)
	
	content := `# Test Book

Some random text
- Not a highlight line
- Almost a highlight (Page: )
- Valid highlight (Page: 123)
Another random line
`
	
	existing := extractExistingHighlights(content)
	
	assert.Equal(1, len(existing), "Should extract only 1 valid highlight")
	assert.True(existing["Valid highlight|123"], "Should find the valid highlight")
}