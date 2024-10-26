package model

import (
	"app/app/enum"

	"github.com/uptrace/bun"
)

type Invitation struct {
	bun.BaseModel `bun:"table:invitations"`

	ID               string                `bun:",default:gen_random_uuid(),pk" json:"id"`
	TeamID           string                `bun:"team_id,notnull" json:"team_id"`       // FK ใช้ชื่อปกติ
	InvitedBy        string                `bun:"invited_by,notnull" json:"invited_by"` // FK ใช้ชื่อปกติ
	InvitedUserEmail string                `bun:"invited_user_email,notnull" json:"invited_user_email"`
	Status           enum.InvitationStatus `bun:"status,notnull" json:"status"`

	CreateUpdateUnixTimestamp
	SoftDelete
}
