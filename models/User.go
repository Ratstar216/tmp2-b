package models

type User struct {
	ID       string `gorm:"primary_key"`
	User_id  string `json:"user_id"`
	Username string `json:"username"`
}
