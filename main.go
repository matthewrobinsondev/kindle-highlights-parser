package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type Highlight struct {
	Title    string
	Author   string
	Page     string
	Location string
	Date     string
	Text     string
}

func main() {
	books, err := ParseClippings("My Clippings.txt")

	if err != nil {
		fmt.Errorf("Unexpected error: %w", err)
		os.Exit(1)
	}

	for title, highlights := range books {
		// Check to see if file exists already
		// Create file for title
		// Format highlights into MD and write to file
	}
}

func ParseClippings(filename string) (map[string][]Highlight, error) {
	file, err := os.Open(filename)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	var lines []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	highlights := make(map[string][]Highlight)
	var currentHighlight Highlight

	title := regexp.MustCompile(`^(.*) \((.*)\)$`)
	metaData := regexp.MustCompile(`Your Highlight.*page ([0-9]+) .*location ([0-9-]+) \| Added on (.*)`)

	for i := 0; i < len(lines); i++ {
		line := lines[i]
		if title.MatchString(line) {
			match := title.FindStringSubmatch(line)

			currentHighlight.Title = match[1]
			currentHighlight.Author = match[2]

			continue
		}

		if metaData.MatchString(line) {
			match := metaData.FindStringSubmatch(line)

			currentHighlight.Page = match[1]
			currentHighlight.Location = match[2]
			currentHighlight.Date = match[3]

			continue
		}

		if len(line) > 0 && !strings.HasPrefix(line, "==========") {
			currentHighlight.Text = line
			highlights[currentHighlight.Title] = append(highlights[currentHighlight.Title], currentHighlight)

			currentHighlight = Highlight{}
		}
	}

	return highlights, nil
}

func WriteMarkdown(highlights map[string][]Highlight) error {

	return nil
}
