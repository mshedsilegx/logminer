package main

// logminer.go contains the core logic for incremental log analysis.
// It manages the process of opening logs, seeking to previous offsets,
// matching regular expressions, and updating the state.

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
)

// LogMiner encapsulates the log parsing operation and its configuration.
type LogMiner struct {
	LogFile   string // Path to the log file to be analyzed
	RegexStr  string // Regular expression pattern for searching
	StateFile string // Path to the file where the search state is stored
	state     State  // Internal state tracking progress
}

// NewLogMiner initializes and returns a new LogMiner instance.
func NewLogMiner(logFile, regexStr, stateFile string) *LogMiner {
	return &LogMiner{
		LogFile:   logFile,
		RegexStr:  regexStr,
		StateFile: stateFile,
	}
}

// Search performs an incremental log search starting from the last saved offset.
// It returns true if the regex was found in any new log lines since the last run.
func (lm *LogMiner) Search() (bool, error) {
	var err error
	// Load the previous state from disk.
	lm.state, err = loadState(lm.StateFile)
	if err != nil {
		return false, err
	}

	// Reset state if the log file name has changed, assuming a new file or rotation.
	if lm.state.Filename != "" && lm.state.Filename != lm.LogFile {
		lm.state = State{}
	}

	// #nosec G304 - filepath.Clean used for mitigation, but dynamic paths are required.
	f, err := os.Open(filepath.Clean(lm.LogFile))
	if err != nil {
		return false, fmt.Errorf("error opening log file: %w", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "error closing log file: %v\n", err)
		}
	}()

	// Seek to the last processed position.
	if _, err := f.Seek(lm.state.Offset, io.SeekStart); err != nil {
		return false, fmt.Errorf("error seeking log file: %w", err)
	}

	// Pre-compile the regex for performance during line-by-line matching.
	re, err := regexp.Compile(lm.RegexStr)
	if err != nil {
		return false, fmt.Errorf("error compiling regex: %w", err)
	}

	// Use a buffered reader for efficient line-by-line reading.
	reader := bufio.NewReader(f)
	found := false
	offset := lm.state.Offset

	for {
		line, err := reader.ReadString('\n')
		lineLen := int64(len(line))

		// Process the line even if it's partial or EOF was reached.
		if lineLen > 0 {
			if re.MatchString(line) {
				found = true
			}
			offset += lineLen
		}

		if err != nil {
			if err == io.EOF {
				break
			}
			return false, fmt.Errorf("error reading log file: %w", err)
		}
	}

	// Update and persist the state for the next run.
	lm.state.Offset = offset
	lm.state.Filename = lm.LogFile
	if err := saveState(lm.StateFile, lm.state); err != nil {
		return false, err
	}

	return found, nil
}
