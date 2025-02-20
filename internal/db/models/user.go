package models

import (
	"time"

	"github.com/google/uuid"
)

const UserName = "employee"

type User struct {
	Id        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey;"`
	Username  string    `gorm:"type:varchar(50);unique;not null"`
	FirstName string    `gorm:"type:varchar(50)"`
	LastName  string    `gorm:"type:varchar(50)"`
	CreatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
}

func (User) TableName() string {
	return UserName
}
