package models

import "fmt"

type User struct {
	ID       int64  `db:"id" json:"id"`
	Username string `db:"username" json:"username"`
	Password string `db:"password" json:"password"`
}

func (u *User) CacheKey() string {
	return fmt.Sprintf("user_%d", u.ID)
}

type RoundPlayer struct {
	UserID int64  `db:"user_id" json:"user_id"`
	IsSpy  bool   `db:"is_spy"  json:"is_spy"`
}


func (u *User) TableName() string {
	return "users"
}
