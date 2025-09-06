package models

import (
	"fmt"
	"time"
)

type Round struct {
	ID        int64     `db:"id" json:"id"`
	GameID    int64     `db:"game_id" json:"game_id"`
	Status    string    `db:"status" json:"status"`   // maps to game_status enum
	Word      string    `db:"word" json:"word"`
	SpyWord   string    `db:"spy_word" json:"spy_word"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

func (r *Round) CacheKey() string {
	return fmt.Sprintf("round_%d", r.ID)
}
