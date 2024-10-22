package model

import (
	"github.com/uptrace/bun"
)

type Team struct {
	bun.BaseModel `bun:"table:teams"`

	ID        int64  `bun:",pk,autoincrement" json:"id"`         // ใช้ ID สำหรับ Primary Key
	TeamName  string `bun:"team_name,notnull" json:"team_name"`
	CreatedBy int64  `bun:"created_by,notnull" json:"created_by"` // FK ใช้ชื่อปกติ

	CreateUpdateUnixTimestamp
	SoftDelete
}
