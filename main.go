package main

// main.go is the entry point for the Log Incremental Miner.
// It handles command-line arguments, initializes the miner, and manages the search process.

import (
	"flag"
	"fmt"
	"os"
)

var version string

func main() {
	// Define and parse command-line flags.
	logFile := flag.String("log", "", "Path to the log file to be analyzed")
	regexStr := flag.String("regex", "", "Regular expression pattern to search for in the log file")
	stateFile := flag.String("state", "logminer.state", "Path to the file where the miner's state will be persisted")
	versionFlag := flag.Bool("version", false, "Display the current version of the tool")

	flag.Parse()

	// Handle the version flag: print version and exit.
	if *versionFlag {
		fmt.Printf("Log Incremental Miner - Version: %s\n", version)
		os.Exit(0)
	}

	// Validate required flags: log file and regex must be provided.
	if *logFile == "" || *regexStr == "" {
		fmt.Println("Usage: logminer -log <log_file> -regex <regex_string> [-state <state_file>]")
		return
	}

	// Initialize the LogMiner with the provided configuration.
	miner := NewLogMiner(*logFile, *regexStr, *stateFile)

	// Execute the search operation.
	found, err := miner.Search()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error during search: %v\n", err)
		os.Exit(1)
	}

	// Output the result of the search (true if matches were found, false otherwise).
	fmt.Println(found)
}
