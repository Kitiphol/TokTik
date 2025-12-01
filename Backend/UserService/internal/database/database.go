package database

import (
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "log"

    "UserService/internal/config"
    "UserService/internal/entity"
)

var DB *gorm.DB

func InitDB() {
    cfg := config.Load()

    db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }

    db.AutoMigrate(&entity.User{}, &entity.Video{}, &entity.Comment{}, &entity.Like{}, &entity.Notification{})

    if err != nil {
        log.Fatal("Failed to migrate database:", err)
    }

    DB = db
    log.Println("Database connection established and migrated successfully")
}
