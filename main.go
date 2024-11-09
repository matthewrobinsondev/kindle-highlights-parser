package main

import (
	"bufio"
	"fmt"
	"log"
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
		fmt.Printf("Unexpected error: %v\n", err)
		os.Exit(1)
	}

	for title, clippings := range books {
		// TODO: phase one implementation, ideally we refactor to do sepeate logic for editing existing files vs creating new ones
		file, err := os.OpenFile(title, os.O_RDWR|os.O_CREATE, 0666)

		if err != nil {
			log.Fatalf("Unexpected error opening file: %v", err)
		}

		markdown := FormatMarkdownForFile(clippings)

		file.WriteString(markdown)
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

	return createHighlights(lines), nil
}

func createHighlights(lines []string) map[string][]Highlight {
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

	return highlights
}

// TODO: one to come back to can I esentially format further to be by paragraphs using page numbers?
func FormatMarkdownForFile(highlights []Highlight) string {
	var markdown strings.Builder

	for _, highlight := range highlights {
		markdown.WriteString(fmt.Sprintf("- %s (Page: %s)\n", highlight.Text, highlight.Page))
	}

	return markdown.String()
}

//
// func WriteMetaData(title string, file *os.File) error {
// 	return nil
// }
