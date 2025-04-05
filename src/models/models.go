package models

import (
    "time"
)

type User struct {
    ID           uint      `gorm:"primaryKey"`
    Username     string    `gorm:"unique;not null"`
    Email        string    `gorm:"unique;not null"`
    PasswordHash string    `gorm:"not null"`
    CreatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP"`
    Videos       []Video   `gorm:"foreignKey:UserID"`
}

type Video struct {
    ID          uint      `gorm:"primaryKey"`
    Title       string    `gorm:"not null"`
    Description string
    FilePath    string    `gorm:"not null"`
    Thumbnail   string
    UserID      uint
    Status      string    `gorm:"default:pending"`
    CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

type VideoProcessing struct {
    ID          uint      `gorm:"primaryKey"`
    VideoID     uint
    KafkaTopic  string    `gorm:"not null"`
    Partition   int       `gorm:"not null"`
    KafkaOffset int64
    Status      string    `gorm:"default:queued"`
    CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}