package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Notification struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null"`
	VideoID   uuid.UUID `gorm:"type:uuid" json:"videoID"`
	Message   string    `gorm:"type:text;not null" json:"message"`
	Read      bool      `gorm:"default:false" json:"read"`
	CreatedAt time.Time `json:"createdAt"`
}



func (u *Notification) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	return
}