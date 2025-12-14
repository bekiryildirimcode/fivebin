package main

import (
	"bytes"
	_ "embed" // Required for go:embed directive
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"
)

//go:embed resources/data.json
var embeddedData []byte

// Dictionary handles memory-efficient operations on the JSON dictionary file
type Dictionary struct {
	filename    string
	totalWords  int
	lastIndex   int
	useEmbedded bool
	dataReader  io.Reader
}

// NewDictionary creates a new Dictionary instance
// If filename is empty, uses embedded data.json
func NewDictionary(filename string) (*Dictionary, error) {
	dict := &Dictionary{
		lastIndex: -1, // -1 means no word has been shown yet
	}

	// Use embedded data if no filename provided or if filename is "embedded"
	if filename == "" || filename == "embedded" {
		dict.useEmbedded = true
		dict.dataReader = bytes.NewReader(embeddedData)
		dict.filename = "embedded"
	} else {
		dict.useEmbedded = false
		dict.filename = filename
	}

	// Count total words on initialization
	count, err := dict.countWords()
	if err != nil {
		return nil, fmt.Errorf("failed to count words: %w", err)
	}

	dict.totalWords = count
	return dict, nil
}

// getReader returns a reader for the dictionary data (file or embedded)
func (d *Dictionary) getReader() (io.Reader, error) {
	if d.useEmbedded {
		// Return a new reader from embedded data
		return bytes.NewReader(embeddedData), nil
	}
	file, err := os.Open(d.filename)
	if err != nil {
		return nil, err
	}
	return file, nil
}

// countWords streams through the JSON file and counts total entries
// This is a one-time operation that uses minimal memory
func (d *Dictionary) countWords() (int, error) {
	reader, err := d.getReader()
	if err != nil {
		return 0, err
	}

	// Close file if it's a file reader
	if file, ok := reader.(*os.File); ok {
		defer file.Close()
	}

	decoder := json.NewDecoder(reader)

	// Read the opening bracket '['
	if _, err := decoder.Token(); err != nil {
		return 0, fmt.Errorf("expected array start: %w", err)
	}

	count := 0
	// Stream through each object without decoding its full content
	for decoder.More() {
		// Skip the entire object by decoding it into a temporary struct
		// This is more reliable than using RawMessage
		var item Root
		if err := decoder.Decode(&item); err != nil {
			return 0, fmt.Errorf("failed to decode entry: %w", err)
		}
		count++
	}

	return count, nil
}

// GetWordByIndex streams to a specific index and returns that word
// Only one word is loaded into memory at a time
func (d *Dictionary) GetWordByIndex(index int) (*Root, error) {
	if index < 0 || index >= d.totalWords {
		return nil, fmt.Errorf("index %d out of range [0, %d)", index, d.totalWords)
	}

	reader, err := d.getReader()
	if err != nil {
		return nil, err
	}

	// Close file if it's a file reader
	if file, ok := reader.(*os.File); ok {
		defer file.Close()
	}

	decoder := json.NewDecoder(reader)

	// Read the opening bracket '['
	if _, err := decoder.Token(); err != nil {
		return nil, fmt.Errorf("expected array start: %w", err)
	}

	currentIndex := 0
	for decoder.More() {
		var item Root
		if err := decoder.Decode(&item); err != nil {
			return nil, fmt.Errorf("failed to decode entry at index %d: %w", currentIndex, err)
		}

		if currentIndex == index {
			return &item, nil
		}
		currentIndex++
	}

	return nil, fmt.Errorf("word at index %d not found", index)
}

// GetRandomWord returns a random word that hasn't been shown in the previous call
// Uses streaming to avoid loading the entire file into memory
func (d *Dictionary) GetRandomWord() (*Root, error) {
	if d.totalWords == 0 {
		return nil, fmt.Errorf("dictionary is empty")
	}

	// Generate a random index, ensuring it's different from the last one
	var randomIndex int
	for {
		randomIndex = rand.Intn(d.totalWords)
		// If we only have one word, we can't avoid repetition
		if d.totalWords == 1 {
			break
		}
		// Ensure we don't select the same word as last time
		if randomIndex != d.lastIndex {
			break
		}
	}

	d.lastIndex = randomIndex

	// Stream to the random index and return that word
	return d.GetWordByIndex(randomIndex)
}

// FindWord searches for a word by name (case-insensitive)
// Uses streaming to avoid loading the entire file into memory
func (d *Dictionary) FindWord(searchWord string) (*Root, error) {
	if searchWord == "" {
		return nil, fmt.Errorf("search word cannot be empty")
	}

	// Normalize search word to lowercase for case-insensitive search
	searchWordLower := strings.ToLower(strings.TrimSpace(searchWord))

	reader, err := d.getReader()
	if err != nil {
		return nil, err
	}

	// Close file if it's a file reader
	if file, ok := reader.(*os.File); ok {
		defer file.Close()
	}

	decoder := json.NewDecoder(reader)

	// Read the opening bracket '['
	if _, err := decoder.Token(); err != nil {
		return nil, fmt.Errorf("expected array start: %w", err)
	}

	// Stream through entries to find matching word
	for decoder.More() {
		var item Root
		if err := decoder.Decode(&item); err != nil {
			return nil, fmt.Errorf("failed to decode entry: %w", err)
		}

		// Case-insensitive comparison
		wordLower := strings.ToLower(strings.TrimSpace(item.Value.Word))
		if wordLower == searchWordLower {
			return &item, nil
		}
	}

	return nil, fmt.Errorf("word '%s' not found", searchWord)
}

// GetTotalWords returns the total number of words in the dictionary
func (d *Dictionary) GetTotalWords() int {
	return d.totalWords
}

// init initializes the random seed
func init() {
	rand.Seed(time.Now().UnixNano())
}
