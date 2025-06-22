package exporter

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/matthewrobinsdev/kindle-notes-parser/internal/config"
	"github.com/matthewrobinsdev/kindle-notes-parser/pkg/models"
)

const (
	dirPermissions    = 0755
	filePermissions   = 0644
	markdownExtension = ".md"
	highlightPrefix   = "- "
	pageFormat        = " (Page: %s)"
	headerFormat      = "# %s\n\n"
	keySeparator      = "|"

	// Regex pattern for parsing existing highlights
	highlightPattern = `^- (.+) \(Page: (\d+)\)$`
)

type FileSystem interface {
	Stat(name string) (os.FileInfo, error)
	ReadFile(filename string) ([]byte, error)
	WriteFile(filename string, data []byte, perm os.FileMode) error
	MkdirAll(path string, perm os.FileMode) error
}

type OSFileSystem struct{}

func (fs OSFileSystem) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

func (fs OSFileSystem) ReadFile(filename string) ([]byte, error) {
	return os.ReadFile(filename)
}

func (fs OSFileSystem) WriteFile(filename string, data []byte, perm os.FileMode) error {
	return os.WriteFile(filename, data, perm)
}

func (fs OSFileSystem) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

type Service struct {
	config      *config.Config
	fs          FileSystem
	highlightRe *regexp.Regexp
}

func New(cfg *config.Config) *Service {
	return NewWithFileSystem(cfg, OSFileSystem{})
}

func NewWithFileSystem(cfg *config.Config, fs FileSystem) *Service {
	return &Service{
		config:      cfg,
		fs:          fs,
		highlightRe: regexp.MustCompile(highlightPattern),
	}
}

func (s *Service) ExportHighlights(bookHighlights map[string][]models.Highlight) ([]models.ExportResult, error) {
	if len(bookHighlights) == 0 {
		return []models.ExportResult{}, nil
	}

	results := make([]models.ExportResult, 0, len(bookHighlights))

	for title, highlights := range bookHighlights {
		if len(highlights) == 0 {
			continue // Skip books with no highlights
		}

		result, err := s.exportBookHighlights(title, highlights)
		if err != nil {
			return results, fmt.Errorf("exporting highlights for %q: %w", title, err)
		}
		results = append(results, result)
	}

	return results, nil
}

func (s *Service) exportBookHighlights(title string, highlights []models.Highlight) (models.ExportResult, error) {
	if title == "" {
		return models.ExportResult{}, fmt.Errorf("book title cannot be empty")
	}

	filename := s.buildFilePath(title)

	if err := s.ensureDirectoryExists(filename); err != nil {
		return models.ExportResult{}, fmt.Errorf("creating directory: %w", err)
	}

	existingContent, existingHighlights, err := s.loadExistingFile(filename, title)
	if err != nil {
		return models.ExportResult{}, fmt.Errorf("loading existing file: %w", err)
	}

	newHighlights, skippedCount := s.filterDuplicates(highlights, existingHighlights)

	if len(newHighlights) > 0 {
		s.appendHighlights(existingContent, newHighlights)

		if err := s.writeFile(filename, existingContent.String()); err != nil {
			return models.ExportResult{}, fmt.Errorf("writing file: %w", err)
		}
	}

	return models.ExportResult{
		BookTitle:    title,
		NewCount:     len(newHighlights),
		SkippedCount: skippedCount,
		TotalCount:   len(highlights),
	}, nil
}

func (s *Service) buildFilePath(title string) string {
	sanitizedTitle := s.sanitizeFilename(title)
	return filepath.Join(s.config.HomeDir, s.config.NotesDirectory, sanitizedTitle+markdownExtension)
}

func (s *Service) sanitizeFilename(filename string) string {
	// Replace common problematic characters
	replacements := map[string]string{
		"/":  "-",
		"\\": "-",
		":":  "-",
		"*":  "",
		"?":  "",
		"\"": "",
		"<":  "",
		">":  "",
		"|":  "-",
	}

	result := filename
	for old, new := range replacements {
		result = strings.ReplaceAll(result, old, new)
	}

	result = strings.TrimSpace(result)
	if len(result) > 200 {
		result = result[:200]
	}

	return result
}

func (s *Service) ensureDirectoryExists(filename string) error {
	dir := filepath.Dir(filename)
	return s.fs.MkdirAll(dir, dirPermissions)
}

func (s *Service) loadExistingFile(filename, title string) (*strings.Builder, map[string]bool, error) {
	content := &strings.Builder{}
	var existingHighlights map[string]bool

	if _, err := s.fs.Stat(filename); err == nil {
		data, err := s.fs.ReadFile(filename)
		if err != nil {
			return nil, nil, fmt.Errorf("reading existing file: %w", err)
		}

		content.Write(data)
		existingHighlights = s.extractExistingHighlights(string(data))

		return content, existingHighlights, nil
	}

	content.WriteString(fmt.Sprintf(headerFormat, title))
	existingHighlights = make(map[string]bool)

	return content, existingHighlights, nil
}

func (s *Service) filterDuplicates(highlights []models.Highlight, existing map[string]bool) ([]models.Highlight, int) {
	newHighlights := make([]models.Highlight, 0, len(highlights))
	skippedCount := 0

	for _, highlight := range highlights {
		key := s.createHighlightKey(highlight)
		if !existing[key] {
			newHighlights = append(newHighlights, highlight)
		} else {
			skippedCount++
		}
	}

	return newHighlights, skippedCount
}

func (s *Service) appendHighlights(content *strings.Builder, highlights []models.Highlight) {
	for _, highlight := range highlights {
		content.WriteString(highlightPrefix)
		content.WriteString(highlight.Text)
		content.WriteString(fmt.Sprintf(pageFormat, highlight.Page))
		content.WriteString("\n")
	}
}

func (s *Service) writeFile(filename, content string) error {
	return s.fs.WriteFile(filename, []byte(content), filePermissions)
}

func (s *Service) extractExistingHighlights(content string) map[string]bool {
	highlights := make(map[string]bool)

	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if matches := s.highlightRe.FindStringSubmatch(line); len(matches) == 3 {
			text := matches[1]
			page := matches[2]
			key := s.buildHighlightKey(text, page)
			highlights[key] = true
		}
	}

	return highlights
}

func (s *Service) createHighlightKey(highlight models.Highlight) string {
	return s.buildHighlightKey(highlight.Text, highlight.Page)
}

func (s *Service) buildHighlightKey(text, page string) string {
	return text + keySeparator + page
}

