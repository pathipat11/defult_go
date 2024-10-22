package model

import (
	"app/app/enum"

	"github.com/uptrace/bun"
)

type TeamMember struct {
	bun.BaseModel `bun:"table:team_members"`

	TeamID     int64         `bun:"team_id,notnull" json:"team_id"` // FK ใช้ชื่อปกติ
	UserID     int64         `bun:"user_id,notnull" json:"user_id"` // FK ใช้ชื่อปกติ
	RoleInTeam enum.CrudRole `bun:"role_in_team,notnull" json:"role_in_team"`

	CreateUpdateUnixTimestamp
	SoftDelete
}
