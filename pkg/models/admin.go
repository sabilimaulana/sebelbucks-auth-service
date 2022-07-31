package models

type Admin struct {
	UserUUID string `json:"uuid"`
	User     User   `json:"user" gorm:"foreignKey:UserUUID"`
}
