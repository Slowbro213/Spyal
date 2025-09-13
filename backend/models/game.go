package models

import (
	"time"
)

type Game struct {
	ID        int64     `db:"id" json:"id"`
	HostID    int64     `db:"host_id" json:"host_id"`
	RoomID    string    `db:"room_id" json:"room_id"`
	Name      string    `db:"name" json:"name"`
	SpyNumber int       `db:"spy_number" json:"spy_number"`
	MaxPlayers int      `db:"max_players" json:"max_players"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	Private   bool      `db:"private" json:"private"`
	Active    bool      `db:"active" json:"active"`
}

func (g *Game) TableName() string {
	return "games"
}
