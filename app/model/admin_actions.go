package model

import "github.com/uptrace/bun"

type AdminAction struct {
	bun.BaseModel `bun:"table:admin_actions"`

	ID      int64  `bun:",pk,autoincrement" json:"id"`      // ใช้ ID สำหรับ Primary Key
	AdminID int64  `bun:"admin_id,notnull" json:"admin_id"` // FK ใช้ชื่อปกติ
	Action  string `bun:"action,notnull" json:"action"`

	CreateUpdateUnixTimestamp
	SoftDelete
}
