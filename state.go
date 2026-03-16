package main

// This file handles the persistence of the log miner's state.
// It allows the application to resume from where it last stopped by saving
// and loading the current file offset and name.

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// State represents the saved state of the log miner.
// It stores the current file position (offset) and the name of the file
// being processed to ensure continuity across restarts.
type State struct {
	Offset   int64  `json:"offset"`
	Filename string `json:"filename"`
}

// loadState reads the state from a file.
// It checks if the state file exists, opens it, and decodes the JSON content
// into a State struct. If the file doesn't exist, it returns an empty State.
func loadState(stateFile string) (State, error) {
	state := State{}
	// Check if state file exists before attempting to open it.
	if _, err := os.Stat(stateFile); err == nil {
		// #nosec G304 - filepath.Clean is used to mitigate path traversal,
		// but the variable path is required for dynamic state file locations.
		f, err := os.Open(filepath.Clean(stateFile))
		if err != nil {
			return state, fmt.Errorf("error opening state file: %w", err)
		}
		// Ensure file is closed after reading.
		defer func() {
			if err := f.Close(); err != nil {
				fmt.Printf("error closing file: %v\n", err)
			}
		}()

		// Decode the JSON-encoded state from the file.
		decoder := json.NewDecoder(f)
		err = decoder.Decode(&state)
		if err != nil {
			return state, fmt.Errorf("error decoding state file: %w", err)
		}
	}
	return state, nil
}

// saveState writes the state to a file.
// It creates or truncates the specified file and encodes the State struct
// as JSON into it.
func saveState(stateFile string, state State) error {
	// #nosec G304 - filepath.Clean is used to mitigate path traversal,
	// but the variable path is required for dynamic state file locations.
	f, err := os.Create(filepath.Clean(stateFile))
	if err != nil {
		return fmt.Errorf("error creating state file: %w", err)
	}
	// Ensure file is closed after writing.
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Printf("error closing file: %v\n", err)
		}
	}()

	// Encode the State struct as JSON into the file.
	encoder := json.NewEncoder(f)
	err = encoder.Encode(state)
	if err != nil {
		return fmt.Errorf("error encoding state file: %w", err)
	}
	return nil
}
