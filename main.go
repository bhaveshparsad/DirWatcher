package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
)

var stopChannel chan bool

func main() {
	dbURL := "user=yourusername password=yourpassword dbname=yourdbname sslmode=disable"
	initDB(dbURL)

	config := &Configuration{
		Directory:      "C:/Users/Bhavesh Prasad/Documents",
		TimeInterval:   5 * time.Second,
		MagicString:    "magic",
		TaskInProgress: false,
	}

	go startBackgroundTask(config)

	// Handle interrupt signals
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig

		log.Println("Received interrupt signal. Stopping background task...")
		stopBackgroundTask()
		// Wait for the background task to finish gracefully
		<-stopChannel
		log.Println("Background task stopped. Exiting...")
		os.Exit(0)
	}()

	// Handle interrupt signals
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig

		log.Println("Received interrupt signal. Stopping background task...")
		stopBackgroundTask()
		// Wait for the background task to finish gracefully
		<-stopChannel
		log.Println("Background task stopped. Exiting...")
		os.Exit(0)
	}()

	startAPIServer(config)
}

func startBackgroundTask(config *Configuration) {
	stopChannel = make(chan bool)

	go func() {
		for {
			select {
			case <-time.After(config.TimeInterval):
				if !config.TaskInProgress {
					go runBackgroundTask(config)
				}
			case <-stopChannel:
				return
			}
		}
	}()
}

func runBackgroundTask(config *Configuration) {
	configMutex.Lock()
	config.TaskInProgress = true
	configMutex.Unlock()

	log.Println("Background task started.")

	startTime := time.Now()

	filesAdded, filesDeleted, magicStringCnt, err := monitorDirectory(config.Directory, config.MagicString)
	if err != nil {
		log.Printf("Error monitoring directory: %v", err)
	}

	endTime := time.Now()
	runtime := endTime.Sub(startTime)

	taskResult := TaskResult{
		StartTime:      startTime,
		EndTime:        endTime,
		Runtime:        runtime,
		FilesAdded:     filesAdded,
		FilesDeleted:   filesDeleted,
		MagicStringCnt: magicStringCnt,
		Status:         "success",
	}

	saveTaskResult(taskResult)

	configMutex.Lock()
	config.TaskInProgress = false
	configMutex.Unlock()

	log.Printf("Background task completed. Magic String Count: %d", magicStringCnt)
}

func monitorDirectory(directory, magicString string) ([]string, []string, int, error) {
	var filesAdded []string
	var filesDeleted []string
	var magicStringCnt int

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, nil, 0, err
	}
	defer watcher.Close()

	// Add directory to watcher
	err = watcher.Add(directory)
	if err != nil {
		return nil, nil, 0, err
	}

	// Keep track of processed events to filter out duplicates on Windows
	processedEvents := make(map[string]struct{})

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return filesAdded, filesDeleted, magicStringCnt, nil
			}

			// Filter out duplicate events on Windows
			if _, processed := processedEvents[event.Name]; processed {
				continue
			}
			processedEvents[event.Name] = struct{}{}

			switch {
			case event.Op&fsnotify.Write == fsnotify.Write:
				// File modified
				log.Printf("Modified file: %s", event.Name)
				content, err := readFileContents(event.Name)
				if err != nil {
					log.Printf("Error reading file contents: %v", err)
					continue
				}
				if strings.Contains(content, magicString) {
					filesAdded = append(filesAdded, event.Name)
					magicStringCnt++
				}
			case event.Op&fsnotify.Create == fsnotify.Create:
				// New file created
				log.Printf("New file created: %s", event.Name)
				content, err := readFileContents(event.Name)
				if err != nil {
					log.Printf("Error reading file contents: %v", err)
					continue
				}
				if strings.Contains(content, magicString) {
					filesAdded = append(filesAdded, event.Name)
					magicStringCnt++
				}
			case event.Op&fsnotify.Remove == fsnotify.Remove || event.Op&fsnotify.Rename == fsnotify.Rename:
				// File removed or renamed
				log.Printf("File removed or renamed: %s", event.Name)
				filesDeleted = append(filesDeleted, event.Name)
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return filesAdded, filesDeleted, magicStringCnt, nil
			}
			log.Printf("Error watching directory: %v", err)
		}
	}
}

func readFileContents(filePath string) (string, error) {
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		content, err := os.ReadFile(filePath)
		if err == nil {
			return string(content), nil
		}

		log.Printf("Error reading file contents (Attempt %d/%d): %v", i+1, maxRetries, err)
	}

	log.Printf("Reached maximum retries, unable to read file contents")
	return "", fmt.Errorf("unable to read file contents after %d attempts", maxRetries)
}

func stopBackgroundTask() {
	close(stopChannel)
}
