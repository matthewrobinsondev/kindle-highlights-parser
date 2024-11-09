package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/viper"
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
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()

	if err != nil {
		fmt.Printf("Errored reading config: %v", err)
		os.Exit(1)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error loading home directory: %v", err)
		os.Exit(1)
	}

	books, err := ParseClippings("My Clippings.txt")
	if err != nil {
		fmt.Printf("Unexpected error: %v\n", err)
		os.Exit(1)
	}

	if !viper.IsSet("notes_directory") {
		fmt.Println("Please add notes_directory to your config")
	}

	notesDirectory := viper.GetString("notes_directory")

	// TODO: Revist this once complete to look into speeding up with concurency
	for title, clippings := range books {
		// TODO: phase one implementation, ideally we refactor to do sepeate logic for editing existing files vs creating new ones
		filename := fmt.Sprintf("%s/%s/%s.md", homeDir, notesDirectory, title)
		file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0666)

		if err != nil {
			log.Fatalf("Unexpected error opening file: %v", err)
		}

		markdown := FormatMarkdownForFile(clippings)
		file.WriteString(markdown)
	}

	fmt.Print("Your notes have now been processed.\n")
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

	var utf8BOM = []byte{0xEF, 0xBB, 0xBF}

	metaData := regexp.MustCompile(`Your Highlight.*page ([0-9]+) .*location ([0-9-]+) \| Added on (.*)`)

	for i := 0; i < len(lines); i++ {
		line := lines[i]
		if title.MatchString(line) {
			match := title.FindStringSubmatch(line)

			// Had to add this due to finding some encoded bytes in the clippings file in titles
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
		// TODO: Using bytes is probably the way to go but quicker to do what I already know
		markdown.WriteString(fmt.Sprintf("- %s (Page: %s)\n", highlight.Text, highlight.Page))
	}

	return markdown.String()
}
