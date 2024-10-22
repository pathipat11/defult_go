package model

import (
	"app/app/enum"

	"github.com/uptrace/bun"
)

type Invitation struct {
	bun.BaseModel `bun:"table:invitations"`

	ID               int64  `bun:",pk,autoincrement" json:"id"`         // ใช้ ID สำหรับ Primary Key
	TeamID           int64  `bun:"team_id,notnull" json:"team_id"`      // FK ใช้ชื่อปกติ
	InvitedBy        int64  `bun:"invited_by,notnull" json:"invited_by"` // FK ใช้ชื่อปกติ
	InvitedUserEmail string `bun:"invited_user_email,notnull" json:"invited_user_email"`
	Status           enum.InvitationStatus `bun:"status,notnull" json:"status"`

	CreateUpdateUnixTimestamp
	SoftDelete
}
