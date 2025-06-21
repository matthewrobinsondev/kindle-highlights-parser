# Kindle Highlights Parser

A terminal UI application using Charm for parsing and organizing Kindle highlights into markdown files.

## Features

- **Book Grouping**: Highlights organized by book with expandable sections
- **Selective Export**: Choose individual highlights or select all from a book
- **Bulk Operations**: Select/deselect all highlights with keyboard shortcuts
- **Markdown Export**: Clean markdown files organized by book title

## Installation & Usage

1. **Get your clippings**: Plug in your Kindle and copy `documents/My Clippings.txt`
2. **Place the file**: Put `My Clippings.txt` in this repository directory
3. **Configure**: Create a `config.toml` file:
   ```toml
   notes_directory = "Documents/your-notes-folder"
   ```
4. **Build and run**:
   ```bash
   go build ./cmd/main.go
   ./kindle-highlights-parser
   ```
   Or run directly:
   ```bash
   go run ./cmd/main.go
   ```

## Development

```bash
# Run tests
go test ./...

# Format code
go fmt ./...
```
