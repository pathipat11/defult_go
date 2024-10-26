package model

import (
	"github.com/uptrace/bun"
)

type TeamMember struct {
	bun.BaseModel `bun:"table:team_members"`

	TeamID string `bun:"team_id,notnull" json:"team_id"` // FK ใช้ชื่อปกติ
	UserID string `bun:"user_id,notnull" json:"user_id"` // FK ใช้ชื่อปกติ

	CreateUpdateUnixTimestamp
}
