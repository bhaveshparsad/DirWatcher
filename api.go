package main

import (
	"time"

	"github.com/gin-gonic/gin"
)

var taskResult TaskResult

func startAPIServer(config *Configuration) {
	r := gin.Default()

	r.POST("/task/start", startBackgroundTaskHandler(config))
	r.GET("/config", getConfigHandler(config))
	r.POST("/config", setConfigHandler(config))
	r.GET("/task/result", getTaskResultHandler)

	r.Run(":8080")
}

func startBackgroundTaskHandler(config *Configuration) gin.HandlerFunc {
	return func(c *gin.Context) {
		go startBackgroundTask(config)
		c.JSON(200, gin.H{"message": "Background task started successfully"})
	}
}

func getConfigHandler(config *Configuration) gin.HandlerFunc {
	return func(c *gin.Context) {
		configMutex.Lock()
		defer configMutex.Unlock()

		c.JSON(200, gin.H{
			"directory":        config.Directory,
			"time_interval":    config.TimeInterval.String(),
			"magic_string":     config.MagicString,
			"task_in_progress": config.TaskInProgress,
		})
	}
}

func setConfigHandler(config *Configuration) gin.HandlerFunc {
	return func(c *gin.Context) {
		configMutex.Lock()
		defer configMutex.Unlock()

		var newConfig Configuration
		if err := c.BindJSON(&newConfig); err != nil {
			c.JSON(400, gin.H{"error": "Invalid JSON format"})
			return
		}

		config.Directory = newConfig.Directory
		config.TimeInterval, _ = time.ParseDuration(newConfig.TimeInterval.String())
		config.MagicString = newConfig.MagicString
		config.TaskInProgress = newConfig.TaskInProgress

		c.JSON(200, gin.H{"message": "Configuration updated successfully"})
	}
}

func getTaskResultHandler(c *gin.Context) {
	if taskResult.StartTime.IsZero() {
		c.JSON(200, gin.H{"message": "No task result available"})
		return
	}

	c.JSON(200, gin.H{
		"start_time":       taskResult.StartTime,
		"end_time":         taskResult.EndTime,
		"runtime":          taskResult.Runtime.String(),
		"files_added":      taskResult.FilesAdded,
		"files_deleted":    taskResult.FilesDeleted,
		"magic_string_cnt": taskResult.MagicStringCnt,
		"status":           taskResult.Status,
	})
}
