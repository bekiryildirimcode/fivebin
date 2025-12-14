package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// Root represents the top-level structure of each dictionary entry
type Root struct {
	ID    int   `json:"id"`
	Value Value `json:"value"`
}

// Value contains the actual word data
type Value struct {
	Word      string    `json:"word"`
	Href      string    `json:"href"`
	Type      string    `json:"type"`
	Level     string    `json:"level"`
	Phonetics Phonetics `json:"phonetics"`
	Examples  []string  `json:"examples"`
}

// Phonetics contains pronunciation information
type Phonetics struct {
	US string `json:"us"`
	UK string `json:"uk"`
}

func main() {
	var dataPath string

	// Check if user provided a custom dictionary file path
	if len(os.Args) > 1 {
		dataPath = os.Args[1]
		// Check if the provided file exists
		if _, err := os.Stat(dataPath); os.IsNotExist(err) {
			log.Fatalf("Dictionary file not found: %s", dataPath)
		}
	} else {
		// Use embedded data.json from resources directory
		dataPath = "embedded"
		fmt.Println("Using embedded dictionary data...")
	}

	// Initialize dictionary with streaming support
	dict, err := NewDictionary(dataPath)
	if err != nil {
		log.Fatalf("Failed to initialize dictionary: %v", err)
	}

	fmt.Printf("Loaded dictionary with %d words\n", dict.GetTotalWords())

	// Initialize user meanings storage
	// Get the executable directory for user meanings
	exePath, err := os.Executable()
	if err != nil {
		log.Fatalf("Failed to get executable path: %v", err)
	}
	exeDir := filepath.Dir(exePath)

	// Check if we're in an app bundle
	var meaningsPath string
	if filepath.Base(exeDir) == "MacOS" {
		// In app bundle, use Resources directory
		meaningsPath = filepath.Join(filepath.Dir(exeDir), "Resources", "user_meanings.json")
	} else {
		// Not in app bundle, use same directory as executable
		meaningsPath = filepath.Join(exeDir, "user_meanings.json")
	}

	meanings, err := NewUserMeanings(meaningsPath)
	if err != nil {
		log.Fatalf("Failed to initialize user meanings: %v", err)
	}

	// Create and run the GUI
	gui := NewGUI(dict, meanings)
	gui.ShowAndRun()
}
