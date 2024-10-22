package model

import "github.com/uptrace/bun"

type ActivateLog struct {
	bun.BaseModel `bun:"table:activate_log"`

	ID          int64  `bun:",pk,autoincrement" json:"id"`    // ใช้ ID สำหรับ Primary Key
	UserID      int64  `bun:"user_id,notnull" json:"user_id"` // FK ใช้ชื่อปกติ
	Action      string `bun:"action,notnull" json:"action"`
	Description string `bun:"description" json:"description"`

	CreateUpdateUnixTimestamp
	SoftDelete
}
