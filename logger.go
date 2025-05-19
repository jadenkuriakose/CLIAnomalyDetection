package main

import (
	"log"
	"os"
	"sync"
)

var (
	logger *log.Logger
	logMu  sync.Mutex
)

func initLogger(filename string) {
	logMu.Lock()
	defer logMu.Unlock()

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	logger = log.New(file, "INFO: ", log.LstdFlags)
}

func logMessage(message string) {
	logMu.Lock()
	defer logMu.Unlock()
	logger.Println(message)
}