package entity

import (
    "github.com/google/uuid"
    "gorm.io/gorm"
)

type Like struct {
    ID      uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
    UserID  uuid.UUID `gorm:"type:uuid;not null;index:idx_user_video,unique" json:"user_id"`
    VideoID uuid.UUID `gorm:"type:uuid;not null;index:idx_user_video,unique" json:"video_id"`
    User    User      `gorm:"foreignKey:UserID" json:"-"`
    Video   Video     `gorm:"foreignKey:VideoID" json:"-"`
}

func (l *Like) BeforeCreate(tx *gorm.DB) (err error) {
    l.ID = uuid.New()
    return
}