package models

import (
	"time"
)

type Round struct {
	ID        int64     `db:"id" json:"id"`
	GameID    int64     `db:"game_id" json:"game_id"`
	Status    string    `db:"status" json:"status"`
	Word      int64     `db:"word_id" json:"word_id"`
	SpyWord   int64     `db:"spy_word_id" json:"spy_word_id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

func (r *Round) TableName() string {
	return "rounds"
}
