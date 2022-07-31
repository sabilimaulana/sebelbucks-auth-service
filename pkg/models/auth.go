package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	UUID      string         `json:"uuid" gorm:"primaryKey"`
	Email     string         `json:"email"`
	Password  string         `json:"password"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime:nano"` // Use unix nano seconds as updating time
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}
