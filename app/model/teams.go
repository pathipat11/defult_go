package model

import (
	"github.com/uptrace/bun"
)

type Team struct {
	bun.BaseModel `bun:"table:teams"`

	ID        string `bun:",default:gen_random_uuid(),pk" json:"id"`
	TeamName  string `bun:"team_name,notnull" json:"team_name"`
	CreatedBy string `bun:"created_by,notnull" json:"created_by"`

	CreateUpdateUnixTimestamp
	SoftDelete
}
