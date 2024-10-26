package model

import "github.com/uptrace/bun"

type Role struct {
	bun.BaseModel `bun:"table:roles"`

	ID          int64  `bun:",pk,autoincrement" json:"id"` // ใช้ ID สำหรับ Primary Key
	RoleName    string `bun:"role_name,notnull" json:"role_name"`
	Description string `bun:"description" json:"description"`

	CreateUpdateUnixTimestamp
	SoftDelete
}
