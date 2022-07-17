package models

type User struct {
	UUID     string `json:"uuid" gorm:"primaryKey"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
