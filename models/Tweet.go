package models

type Tweet struct {
	ID         string `json:"id" gorm:"primary_key"`
	UserID     string `json:"user_id"`
	ParentID   string `json:"parent_id"`
	Username   string `json:"username"`
	Likes      int    `json:"likes"`
	Content    string `json:"content"`
	Created_at string `json:"created_at"`
	// Replies    []Tweet `json:"replies"`
}

// type Tweet struct {
// 	id         string `gorm:"primary_key"`
// 	user_id    string
// 	parent_id  string
// 	username   string
// 	likes      int
// 	content    string
// 	created_at string
// }
