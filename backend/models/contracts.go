package models

type Model interface {
	any
	TableName() string
}
