package model

import "github.com/uptrace/bun"

type Notification struct {
	bun.BaseModel `bun:"table:notifications"`

	ID        int64  `bun:",pk,autoincrement" json:"id"`          // ใช้ ID สำหรับ Primary Key
	UserID    int64  `bun:"user_id,notnull" json:"user_id"`       // FK ใช้ชื่อปกติ
	Message   string `bun:"message" json:"message"`
	IsRead    bool   `bun:"is_read" json:"is_read"`

	CreateUpdateUnixTimestamp
	SoftDelete
}
