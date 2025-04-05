package models

import (
    "time"
)

type User struct {
    ID           uint      `gorm:"primaryKey;autoIncrement"`
    Username     string    `gorm:"type:varchar(50);unique;not null"`
    Email        string    `gorm:"type:varchar(100);unique;not null"`
    PasswordHash string    `gorm:"type:text;not null"`
    CreatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP"`
    Videos       []Video   `gorm:"foreignKey:UserID"`
}

type Video struct {
    ID          uint      `gorm:"primaryKey;autoIncrement"`
    Title       string    `gorm:"type:varchar(255);not null"`
    Description string    `gorm:"type:text"`
    FilePath    string    `gorm:"type:text;not null"`
    Thumbnail   string    `gorm:"type:text"`
    UserID      uint      `gorm:"not null"`
    Status      string    `gorm:"type:varchar(20);default:pending;check:status IN ('pending', 'processing', 'ready')"`
    CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP"`
    User        User      `gorm:"foreignKey:UserID"`
}

type VideoProcessing struct {
    ID          uint      `gorm:"primaryKey;autoIncrement"`
    VideoID     uint      `gorm:"not null"`
    KafkaTopic  string    `gorm:"type:varchar(255);not null"`
    Partition   int       `gorm:"not null"`
    KafkaOffset int64
    Status      string    `gorm:"type:varchar(20);default:queued;check:status IN ('queued', 'processing', 'completed', 'failed')"`
    CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP"`
    Video       Video     `gorm:"foreignKey:VideoID"`
}