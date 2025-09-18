package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
)

// LogMiner encapsulates the log parsing operation.
type LogMiner struct {
	LogFile   string
	RegexStr  string
	StateFile string
	state     State
}

// NewLogMiner creates a new LogMiner.
func NewLogMiner(logFile, regexStr, stateFile string) *LogMiner {
	return &LogMiner{
		LogFile:   logFile,
		RegexStr:  regexStr,
		StateFile: stateFile,
	}
}

// Search performs the incremental log search.
func (lm *LogMiner) Search() (bool, error) {
	var err error
	lm.state, err = loadState(lm.StateFile)
	if err != nil {
		return false, err
	}

	if lm.state.Filename != "" && lm.state.Filename != lm.LogFile {
		lm.state = State{}
	}

	f, err := os.Open(lm.LogFile)
	if err != nil {
		return false, fmt.Errorf("error opening log file: %w", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Printf("error closing file: %v\n", err)
		}
	}()

	if _, err := f.Seek(lm.state.Offset, io.SeekStart); err != nil {
		return false, fmt.Errorf("error seeking log file: %w", err)
	}

	re, err := regexp.Compile(lm.RegexStr)
	if err != nil {
		return false, fmt.Errorf("error compiling regex: %w", err)
	}

	reader := bufio.NewReader(f)
	found := false
	offset := lm.state.Offset

	for {
		line, err := reader.ReadString('\n')
		lineLen := int64(len(line))

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

	lm.state.Offset = offset
	lm.state.Filename = lm.LogFile
	if err := saveState(lm.StateFile, lm.state); err != nil {
		return false, err
	}

	return found, nil
}
