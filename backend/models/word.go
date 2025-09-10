package models

import "fmt"

type Word struct {
	ID   int64  `db:"id"   json:"id"`
	Word string `db:"word" json:"word"`
}

func (w *Word) CacheKey() string {
	return fmt.Sprintf("word_%d", w.ID)
}
