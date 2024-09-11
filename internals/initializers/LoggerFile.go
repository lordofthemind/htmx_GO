package initializers

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

func SetUpLoggerFile(logFileName string) (*os.File, error) {
	// Ensure the logs directory exists
	if err := os.MkdirAll("logs", 0755); err != nil {
		return nil, fmt.Errorf("failed to create logs directory: %w", err)
	}

	// Get the current date and time for the log file prefix
	currentTime := time.Now().Format("20060102_150405")
	logFileName = fmt.Sprintf("%s_%s", currentTime, logFileName)

	// Construct the full log file path
	logFilePath := filepath.Join("logs", logFileName)

	// Check if the log file path is a directory
	if stat, err := os.Stat(logFilePath); err == nil && stat.IsDir() {
		log.Printf("Error: %s is a directory, not a file", logFilePath)
		return nil, fmt.Errorf("log file path %s is a directory, not a file", logFilePath)
	}

	// Open the log file
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		// If there's an error opening the log file, fallback to stdout only logging
		log.Printf("Error opening log file %s, falling back to stdout only: %v", logFilePath, err)
		return nil, fmt.Errorf("failed to open log file: %w", err)
	} else {
		// Successfully opened the log file
		log.Printf("Logging initialized. Log file: %s", logFilePath)
	}

	// Set up multi-writer to write to stdout and file if possible
	var multiWriter io.Writer
	if logFile != nil {
		multiWriter = io.MultiWriter(os.Stdout, logFile)
	} else {
		multiWriter = os.Stdout
	}

	// Configure log format to include timestamp, log level, and file location
	log.SetOutput(multiWriter)
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)

	return logFile, nil
}
