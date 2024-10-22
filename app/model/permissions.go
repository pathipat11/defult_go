package model

import "github.com/uptrace/bun"

type Permission struct {
	bun.BaseModel `bun:"table:permissions"`

	ID            int64  `bun:",pk,autoincrement" json:"id"`   // ใช้ ID สำหรับ Primary Key
	PermissionName string `bun:"permission_name,notnull" json:"permission_name"`
	Description    string `bun:"description" json:"description"`

	CreateUpdateUnixTimestamp
	SoftDelete
}
