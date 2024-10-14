package main

import (
	"os"
	"testing"
)

func TestParseClippings(t *testing.T) {
	filePath := "testData/Test Clippings.txt"

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Fatalf("Test file not found: %v", err)
	}

	ParseClippings(filePath)
}
