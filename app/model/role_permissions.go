package model

import "github.com/uptrace/bun"

type RolePermission struct {
	bun.BaseModel `bun:"table:role_permissions"`

	RoleID       int64 `bun:"role_id,notnull" json:"role_id"`       // FK ใช้ชื่อปกติ
	PermissionID int64 `bun:"permission_id,notnull" json:"permission_id"` // FK ใช้ชื่อปกติ

	CreateUpdateUnixTimestamp
	SoftDelete
}
