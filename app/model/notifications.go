package model

import "github.com/uptrace/bun"

type Notification struct {
	bun.BaseModel `bun:"table:notifications"`

	ID      string `bun:",default:gen_random_uuid(),pk" json:"id"`
	UserID  string `bun:"user_id,notnull" json:"user_id"` // FK ใช้ชื่อปกติ
	Message string `bun:"message" json:"message"`
	IsRead  bool   `bun:"is_read" json:"is_read"`

	CreateUpdateUnixTimestamp
	SoftDelete
}
