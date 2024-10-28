package model

import (
	"app/app/enum"

	"github.com/uptrace/bun"
)

type SprintPlan struct {
	bun.BaseModel `bun:"table:sprint_plan"`

	ID         string            `bun:",default:gen_random_uuid(),pk" json:"id"`
	SprintName string            `bun:"sprint_name,notnull" json:"sprint_name"`
	TeamID     string            `bun:"team_id,notnull" json:"team_id"`
	StartDate  int64             `bun:"start_date,notnull" json:"start_date"`
	EndDate    int64             `bun:"end_date,notnull" json:"end_date"`
	Status     enum.SprintStatus `bun:"status,notnull" json:"status"`

	CreateUpdateUnixTimestamp
	SoftDelete
}
