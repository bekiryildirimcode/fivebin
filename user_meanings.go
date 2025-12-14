package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

// UserMeanings manages user-saved meanings/translations/notes
// Uses a map with word as key and meaning as value
type UserMeanings struct {
	filename string
	meanings map[string]string // word -> meaning
	mu       sync.RWMutex       // For thread-safe access
}

// NewUserMeanings creates a new UserMeanings instance
// Loads existing meanings from file if it exists
func NewUserMeanings(filename string) (*UserMeanings, error) {
	um := &UserMeanings{
		filename: filename,
		meanings: make(map[string]string),
	}

	// Try to load existing meanings
	if err := um.load(); err != nil {
		// If file doesn't exist, that's okay - start with empty map
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("failed to load user meanings: %w", err)
		}
	}

	return um, nil
}

// load reads user meanings from the JSON file
func (um *UserMeanings) load() error {
	um.mu.Lock()
	defer um.mu.Unlock()

	file, err := os.Open(um.filename)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&um.meanings); err != nil {
		return fmt.Errorf("failed to decode user meanings: %w", err)
	}

	return nil
}

// save writes user meanings to the JSON file
func (um *UserMeanings) save() error {
	um.mu.RLock()
	defer um.mu.RUnlock()

	file, err := os.Create(um.filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Pretty print JSON
	if err := encoder.Encode(um.meanings); err != nil {
		return fmt.Errorf("failed to encode user meanings: %w", err)
	}

	return nil
}

// Get retrieves the meaning for a word
// Returns empty string if word not found
func (um *UserMeanings) Get(word string) string {
	um.mu.RLock()
	defer um.mu.RUnlock()

	return um.meanings[word]
}

// Set saves a meaning for a word
func (um *UserMeanings) Set(word, meaning string) error {
	um.mu.Lock()
	um.meanings[word] = meaning
	um.mu.Unlock()

	// Save to file
	return um.save()
}

// Has checks if a meaning exists for a word
func (um *UserMeanings) Has(word string) bool {
	um.mu.RLock()
	defer um.mu.RUnlock()

	_, exists := um.meanings[word]
	return exists
}

