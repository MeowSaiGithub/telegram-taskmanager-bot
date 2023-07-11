package model

import "time"

type Task struct {
	Id          string    `gorm:"column:id;primaryKey;index"`
	CreatedAt   time.Time `gorm:"created_at"`
	ChatId      int64     `gorm:"column:chat_id;index"`
	Description string    `gorm:"column:description"`
}
