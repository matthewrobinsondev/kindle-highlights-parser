package parser

import (
	"bufio"
	"bytes"
	"os"
	"regexp"
	"strings"

	"github.com/matthewrobinsdev/kindle-notes-parser/pkg/models"
)

func ParseClippings(filename string) ([]models.Highlight, error) {
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

func createHighlights(lines []string) []models.Highlight {
	var highlights []models.Highlight
	var currentHighlight models.Highlight

	title := regexp.MustCompile(`^(.*) \((.*)\)$`)
	var utf8BOM = []byte{0xEF, 0xBB, 0xBF}
	metaData := regexp.MustCompile(`Your Highlight.*page ([0-9]+) .*location ([0-9-]+) \| Added on (.*)`)

	for i := range lines {
		line := lines[i]
		if title.MatchString(line) {
			match := title.FindStringSubmatch(line)

			byteText := []byte(match[1])
			if bytes.HasPrefix(byteText, utf8BOM) {
				match[1] = string(byteText[len(utf8BOM):])
			}

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
			highlights = append(highlights, currentHighlight)
			currentHighlight = models.Highlight{}
		}
	}

	return highlights
}

func GroupHighlightsByBook(highlights []models.Highlight) []models.BookGroup {
	bookMap := make(map[string][]models.Highlight)
	authorMap := make(map[string]string)

	for _, highlight := range highlights {
		bookMap[highlight.Title] = append(bookMap[highlight.Title], highlight)
		authorMap[highlight.Title] = highlight.Author
	}

	var books []models.BookGroup
	for title, highlights := range bookMap {
		books = append(books, models.BookGroup{
			Title:      title,
			Author:     authorMap[title],
			Highlights: highlights,
			Expanded:   false,
		})
	}

	return books
}
