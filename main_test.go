package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const CLIPPINGS_FILE_PATH = "testData/Test Clippings.txt"
const FORMATTED_MARKDOWN_FILE_PATH = "testData/SandwormFormatted.md"

func TestParseClippings(t *testing.T) {
	assert := assert.New(t)

	if _, err := os.Stat(CLIPPINGS_FILE_PATH); os.IsNotExist(err) {
		t.Fatalf("Test file not found: %v", err)
	}

	highlights, err := ParseClippings(CLIPPINGS_FILE_PATH)

	if err != nil {
		t.Errorf("Unexpected error parising clippings")
	}

	assert.Equal(2, len(highlights), "There should be two highlights")
	assert.Equal(2, len(highlights["Sandworm"]), "There should be two highlights for this book")
	assert.Equal(1, len(highlights["Modern Software Engineering"]), "There should be one highlight for this book")
	assert.Equal("Put more simply, a complex system like a digitized civilization is subject to cascading failures, where one thing depends on another, which depends on another thing.", highlights["Sandworm"][0].Text)
	assert.Equal("Greenberg, Andy", highlights["Sandworm"][0].Author)
}

func TestFormatMarkdown(t *testing.T) {
	highlights, err := ParseClippings(CLIPPINGS_FILE_PATH)

	if err != nil {
		t.Errorf("Unexpected error parising clippings")
	}

	result := FormatMarkdownForFile(highlights["Sandworm"])

	file, err := os.ReadFile(FORMATTED_MARKDOWN_FILE_PATH)
	if err != nil {
		t.Fatalf("Formatted markdown example not found: %v", err)
	}

	s := string(file)
	assert := assert.New(t)
	assert.Equal(s, result, "Sandworm markdown should match")
}
