package main

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func initDB(dbURL string) (*gorm.DB, error) {
    db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
    if err != nil {
        log.Printf("Error connecting to the database: %v", err)
        return nil, err
    }

    // AutoMigrate creates or updates the database schema
    err = db.AutoMigrate(&TaskResult{})
    if err != nil {
        log.Printf("Error auto-migrating database schema: %v", err)
        return nil, err
    }

    return db, nil
}

func saveTaskResult(result TaskResult) {
	if err := db.Create(&result).Error; err != nil {
		log.Println("Error inserting task result:", err)
		return
	}
}
