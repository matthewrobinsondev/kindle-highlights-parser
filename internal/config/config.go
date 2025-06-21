package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	NotesDirectory string
	HomeDir        string
}

func Load() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Warning: Could not read config file: %v", err)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Error loading home directory: %v", err)
	}

	notesDirectory := viper.GetString("notes_directory")
	if notesDirectory == "" {
		notesDirectory = "notes"
	}

	return &Config{
		NotesDirectory: notesDirectory,
		HomeDir:        homeDir,
	}
}
