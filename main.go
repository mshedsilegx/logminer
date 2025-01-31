package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
)

type State struct {
	LastLine int    `json:"last_line"`
	Found    bool   `json:"found"`
	Filename string `json:"filename"` // Add filename to state
}

func parseLog(logFile, regexStr, stateFile string) (bool, error) {
	// Load state
	state := State{}
	if _, err := os.Stat(stateFile); err == nil {
		f, err := os.Open(stateFile)
		if err != nil {
			return false, fmt.Errorf("error opening state file: %v", err)
		}
		defer f.Close()
		decoder := json.NewDecoder(f)
		err = decoder.Decode(&state)
		if err != nil {
			return false, fmt.Errorf("error decoding state file: %v", err)
		}

	}
	//Check if filename matches
	if state.Filename != "" && state.Filename != logFile {
		// Reset state if filename has changed
		state = State{}
	}

	// Open log file
	f, err := os.Open(logFile)
	if err != nil {
		return false, fmt.Errorf("error opening log file: %v", err)
	}
	defer f.Close()

	reader := bufio.NewReader(f)
	re, err := regexp.Compile(regexStr)
	if err != nil {
		return false, fmt.Errorf("error compiling regex: %v", err)
	}

	lineCount := 0
	for {
		line, err := reader.ReadString('\n')
		lineCount++

		if err != nil {
			if err == io.EOF {
				break // End of file
			}
			return false, fmt.Errorf("error reading log file: %v", err)
		}

		if lineCount > state.LastLine { // Only process new lines
			if re.MatchString(line) {
				state.Found = true
				state.LastLine = lineCount
				break // Found the pattern, no need to continue
			}
			state.LastLine = lineCount //Update last read line
		}

	}

	// Save state
	state.Filename = logFile //Store the filename
	f, err = os.Create(stateFile)
	if err != nil {
		return false, fmt.Errorf("error creating state file: %v", err)
	}
	defer f.Close()

	encoder := json.NewEncoder(f)
	err = encoder.Encode(state)
	if err != nil {
		return false, fmt.Errorf("error encoding state file: %v", err)
	}

	return state.Found, nil
}

func main() {
	logFile := flag.String("log", "", "Path to the log file")
	regexStr := flag.String("regex", "", "Regular expression to search")
	stateFile := flag.String("state", "logminer.state", "Path to the state file")

	flag.Parse()

	if *logFile == "" || *regexStr == "" {
		fmt.Println("Usage: logminer -log <log_file> -regex <regex_string> [-state <state_file>]")
		return
	}

	found, err := parseLog(*logFile, *regexStr, *stateFile)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	fmt.Println(found)
}
