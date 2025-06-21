package main

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const CLIPPINGS_FILE_PATH = "testData/Test Clippings.txt"
const STRANGE_CLIPPINGS_FILE_PATH = "testData/Strange Title Clippings.txt"
const FORMATTED_MARKDOWN_FILE_PATH = "testData/SandwormFormatted.md"

func TestParseClippings(t *testing.T) {
	assert := assert.New(t)

	if _, err := os.Stat(CLIPPINGS_FILE_PATH); os.IsNotExist(err) {
		t.Fatalf("Test file not found: %v", err)
	}

	highlights, err := parseClippings(CLIPPINGS_FILE_PATH)

	if err != nil {
		t.Errorf("Unexpected error parsing clippings")
	}

	sandwormCount := 0
	modernSoftwareCount := 0

	for _, highlight := range highlights {
		if highlight.Title == "Sandworm" {
			sandwormCount++
		}
		if highlight.Title == "Modern Software Engineering" {
			modernSoftwareCount++
		}
	}

	assert.Equal(2, sandwormCount, "There should be two highlights for Sandworm")
	assert.Equal(1, modernSoftwareCount, "There should be one highlight for Modern Software Engineering")

	var sandwormHighlight Highlight
	for _, highlight := range highlights {
		if highlight.Title == "Sandworm" {
			sandwormHighlight = highlight
			break
		}
	}

	assert.Equal("Put more simply, a complex system like a digitized civilization is subject to cascading failures, where one thing depends on another, which depends on another thing.", sandwormHighlight.Text)
	assert.Equal("Greenberg, Andy", sandwormHighlight.Author)
}

func TestStrangBooknameCanBeParsed(t *testing.T) {
	assert := assert.New(t)

	if _, err := os.Stat(STRANGE_CLIPPINGS_FILE_PATH); os.IsNotExist(err) {
		t.Fatalf("Test file not found: %v", err)
	}

	highlights, err := parseClippings(STRANGE_CLIPPINGS_FILE_PATH)

	if err != nil {
		t.Errorf("Unexpected error parsing clippings")
	}

	sandwormCount := 0
	for _, highlight := range highlights {
		if highlight.Title == "Sandworm" {
			sandwormCount++
		}
	}

	assert.Equal(2, sandwormCount, "There should be two highlights for Sandworm")
}

func TestFormatMarkdown(t *testing.T) {
	highlights, err := parseClippings(CLIPPINGS_FILE_PATH)

	if err != nil {
		t.Errorf("Unexpected error parsing clippings")
	}

	var sandwormHighlights []Highlight
	for _, highlight := range highlights {
		if highlight.Title == "Sandworm" {
			sandwormHighlights = append(sandwormHighlights, highlight)
		}
	}

	result := formatMarkdownForFile(sandwormHighlights)

	file, err := os.ReadFile(FORMATTED_MARKDOWN_FILE_PATH)
	if err != nil {
		t.Fatalf("Formatted markdown example not found: %v", err)
	}

	s := string(file)
	assert := assert.New(t)
	assert.Equal(s, result, "Sandworm markdown should match")
}

func formatMarkdownForFile(highlights []Highlight) string {
	var markdown strings.Builder

	for _, highlight := range highlights {
		markdown.WriteString(fmt.Sprintf("- %s (Page: %s)\n", highlight.Text, highlight.Page))
	}

	return markdown.String()
}
