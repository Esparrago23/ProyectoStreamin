package database

import (
    "fmt"
    "os"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "streaming-service/src/models"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
    dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
        os.Getenv("DB_HOST"),
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASSWORD"),
        os.Getenv("DB_NAME"),
        os.Getenv("DB_PORT"),
    )

    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        panic("Failed to connect to database")
    }

    // Create the tables
    err = db.AutoMigrate(
        &models.User{},
        &models.Video{},
        &models.VideoProcessing{},
    )
    if err != nil {
        panic("Failed to migrate database: " + err.Error())
    }

    DB = db
    return db
}