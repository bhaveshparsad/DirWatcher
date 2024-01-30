package main

import (
	"sync"
	"time"
)

var configMutex sync.Mutex

// Configuration holds the configuration for the DirWatcher application
type Configuration struct {
	Directory      string
	TimeInterval   time.Duration
	MagicString    string
	TaskInProgress bool
}
