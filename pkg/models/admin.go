package models

import (
	"time"

	"gorm.io/gorm"
)

type Admin struct {
	UserUUID  string         `json:"uuid"`
	User      User           `json:"user" gorm:"foreignKey:UserUUID"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime:nano"` // Use unix nano seconds as updating time
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}
