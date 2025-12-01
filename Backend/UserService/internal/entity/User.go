// entity/user.go
package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)


type User struct {
    ID           uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
    Username     string    `gorm:"uniqueIndex;size:100;not null" json:"username"`
    PasswordHash string    `gorm:"size:255;not null" json:"-"`
    Email        string    `gorm:"uniqueIndex;size:100;not null" json:"email"`

    Videos   []Video   `gorm:"foreignKey:UserID" json:"videos,omitempty"`    // One-to-many: User has many Videos
    // Comments []Comment `gorm:"foreignKey:UserID" json:"comments,omitempty"`  // One-to-many: User has many Comments
}

// Auto-generate UUID before creating
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	return
}






