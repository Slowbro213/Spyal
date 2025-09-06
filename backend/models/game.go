package models

import (
	"fmt"
	"time"
)

type Game struct {
	ID        int64     `db:"id" json:"id"`
	HostID    int64     `db:"host_id" json:"host_id"`
	Title     string    `db:"title" json:"title"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	Private   bool      `db:"private" json:"private"`
	Status    string    `db:"status" json:"status"`
}

func (g *Game) CacheKey() string {
	return fmt.Sprintf("game_%d", g.ID)
}
