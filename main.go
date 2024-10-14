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
	ParseClippings("My Clippings.txt")
}

func ParseClippings(filename string) (map[string][]Highlight, error) {
	file, err := os.Open(filename)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	var lines []string

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
		}
	}

	fmt.Printf("%+v\n", highlights)

	return highlights, nil
}
