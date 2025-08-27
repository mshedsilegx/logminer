package main

import (
	"flag"
	"fmt"
	"os"
)

var version string

func main() {
	logFile := flag.String("log", "", "Path to the log file")
	regexStr := flag.String("regex", "", "Regular expression to search")
	stateFile := flag.String("state", "logminer.state", "Path to the state file")
	versionFlag := flag.Bool("version", false, "Display version information")

	flag.Parse()

	// Version
	if *versionFlag {
		fmt.Printf("Log Incremental Miner - Version: %s\n", version)
		os.Exit(0)
	}

	if *logFile == "" || *regexStr == "" {
		fmt.Println("Usage: logminer -log <log_file> -regex <regex_string> [-state <state_file>]")
		return
	}

	miner := NewLogMiner(*logFile, *regexStr, *stateFile)
	found, err := miner.Search()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	fmt.Println(found)
}
