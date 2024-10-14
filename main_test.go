package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseClippings(t *testing.T) {
	assert := assert.New(t)
	filePath := "testData/Test Clippings.txt"

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Fatalf("Test file not found: %v", err)
	}

	highlights, err := ParseClippings(filePath)

	if err != nil {
		t.Errorf("Unexpected error parising clippings")
	}

	assert.Equal(2, len(highlights), "There should be two highlights")
	assert.Equal("Put more simply, a complex system like a digitized civilization is subject to cascading failures, where one thing depends on another, which depends on another thing.", highlights["Sandworm"][0].Text)
	assert.Equal("Greenberg, Andy", highlights["Sandworm"][0].Author)
}
