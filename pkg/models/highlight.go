package models

type Highlight struct {
	Title    string
	Author   string
	Page     string
	Location string
	Date     string
	Text     string
}

type BookGroup struct {
	Title      string
	Author     string
	Highlights []Highlight
	Expanded   bool
}

type ListItem struct {
	IsBook         bool
	BookIndex      int
	HighlightIndex int
}

type ExportResult struct {
	BookTitle    string
	NewCount     int
	SkippedCount int
	TotalCount   int
}
