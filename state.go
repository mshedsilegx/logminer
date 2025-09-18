package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// State represents the saved state of the log miner.
type State struct {
	Offset   int64  `json:"offset"`
	Filename string `json:"filename"`
}

// loadState reads the state from a file.
func loadState(stateFile string) (State, error) {
	state := State{}
	if _, err := os.Stat(stateFile); err == nil {
		f, err := os.Open(stateFile)
		if err != nil {
			return state, fmt.Errorf("error opening state file: %w", err)
		}
		defer func() {
			if err := f.Close(); err != nil {
				fmt.Printf("error closing file: %v\n", err)
			}
		}()
		decoder := json.NewDecoder(f)
		err = decoder.Decode(&state)
		if err != nil {
			return state, fmt.Errorf("error decoding state file: %w", err)
		}
	}
	return state, nil
}

// saveState writes the state to a file.
func saveState(stateFile string, state State) error {
	f, err := os.Create(stateFile)
	if err != nil {
		return fmt.Errorf("error creating state file: %w", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Printf("error closing file: %v\n", err)
		}
	}()

	encoder := json.NewEncoder(f)
	err = encoder.Encode(state)
	if err != nil {
		return fmt.Errorf("error encoding state file: %w", err)
	}
	return nil
}
