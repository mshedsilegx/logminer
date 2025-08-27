# Log Incremental Miner (`logminer`)

## Overview

`logminer` is a command-line utility written in Go that efficiently searches for a regular expression pattern within a log file. It's designed for incremental searching, meaning it keeps track of its progress and will only search new content added to a log file since its last run.

This is particularly useful for monitoring log files for specific events without having to re-scan the entire file each time.

## Objectives

*   **Incremental Searching**: Avoid re-processing the entire log file on every execution. The state (last position read) is saved in a state file.
*   **Efficiency**: By using byte offsets instead of line numbers, `logminer` can quickly seek to the last known position in the file, making it very fast even for very large log files.
*   **Simplicity**: Provide a simple and intuitive command-line interface.

## Command-Line Usage

The tool is run from the command line with the following arguments:

```sh
logminer -log <path_to_log_file> -regex <regex_pattern> [-state <path_to_state_file>]
```

### Arguments:

*   `-log` (required): The path to the log file you want to search.
*   `-regex` (required): The regular expression pattern to search for.
*   `-state` (optional): The path to the state file. If not provided, it defaults to `logminer.state` in the current directory.
*   `-version` (optional): Displays the version information for the tool.

## Examples

### First Run

Imagine you have a log file `app.log` and you want to search for the word "ERROR".

```sh
# Create a sample log file
echo "INFO: Application started." > app.log
echo "DEBUG: Connecting to database." >> app.log
echo "ERROR: Database connection failed." >> app.log

# Run the logminer
logminer -log app.log -regex "ERROR"
```

**Output:**
```
true
```
This command will print `true` because "ERROR" is found. It will also create a `logminer.state` file to save its progress.

### Subsequent Runs

If you run the same command again without any changes to `app.log`, it will start from where it left off and find no new matches.

```sh
# Run the logminer again
logminer -log app.log -regex "ERROR"
```

**Output:**
```
false
```

Now, let's add a new error to the log file.

```sh
# Add a new line to the log file
echo "ERROR: Authentication service timed out." >> app.log

# Run the logminer one more time
logminer -log app.log -regex "ERROR"
```

**Output:**
```
true
```
This time, it outputs `true` again because it found the new "ERROR" line that was added after its last run.

### Advanced Regexp Examples

The `-regex` flag accepts any valid Go regular expression. This allows for more powerful and specific searches.

**Example 1: Searching for errors from a specific service**

Imagine your log lines are formatted like `[SERVICE_NAME] [LOG_LEVEL]: Message`. You only want to find errors from `auth-service`.

```sh
# Sample log lines
# [auth-service] ERROR: Invalid credentials
# [payment-service] INFO: Payment processed
# [auth-service] DEBUG: User lookup
# [payment-service] ERROR: Card declined

logminer -log app.log -regex "\[auth-service\] ERROR"
```
*Note: The `[` and `]` are escaped with backslashes because they are special characters in regular expressions.*

**Example 2: Searching for a pattern with wildcards**

Suppose you want to find any log line that contains an HTTP status code in the 4xx range (e.g., 400, 401, 404).

```sh
# Sample log lines
# Request failed with status code 404
# Request succeeded with status code 200
# Request failed with status code 401

logminer -log app.log -regex "status code 4[0-9]{2}"
```
This regex looks for the literal string "status code 4" followed by any two digits from 0-9.
