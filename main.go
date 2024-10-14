package main

import (
	"bufio"
	"fmt"
	"os"
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

func ParseClippings(filename string) {
	file, err := os.Open(filename)

	if err != nil {
		return
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	var lines []string

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	fmt.Println(lines)
}
