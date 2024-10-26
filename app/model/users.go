package model

import (
	"app/app/enum"

	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users"`

	ID          string      `bun:",default:gen_random_uuid(),pk" json:"id"`
	Username    string      `bun:"username" json:"username"`
	Email       string      `bun:"email" json:"email"`
	Password    string      `bun:"password" json:"password"`
	FirstName   string      `bun:"first_name" json:"first_name"`
	LastName    string      `bun:"last_name" json:"last_name"`
	DisplayName string      `bun:"display_name" json:"display_name"`
	RoleID      int64       `bun:"role_id,notnull" json:"role_id"`
	Status      enum.Status `bun:"status,notnull" json:"status"`

	CreateUpdateUnixTimestamp
	SoftDelete
}
