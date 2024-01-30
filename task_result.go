package main

import (
	"time"
	"gorm.io/gorm"
)

// TaskResult represents the result of a background task run
type TaskResult struct {
	gorm.Model
	StartTime      time.Time
	EndTime        time.Time
	Runtime        time.Duration
	FilesAdded     []string `gorm:"type:text[]"`
	FilesDeleted   []string `gorm:"type:text[]"`
	MagicStringCnt int
	Status         string
}