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
