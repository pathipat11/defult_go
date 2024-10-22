package model

import (
	"app/app/enum"

	"github.com/uptrace/bun"
)

type SprintPlan struct {
	bun.BaseModel `bun:"table:sprint_plan"`

	ID         int64  `bun:",pk,autoincrement" json:"id"` // ใช้ ID สำหรับ Primary Key
	SprintName string `bun:"sprint_name,notnull" json:"sprint_name"`
	TeamID     int64  `bun:"team_id,notnull" json:"team_id"` // FK ใช้ชื่อปกติ
	StartDate  string `bun:"start_date,notnull" json:"start_date"`
	EndDate    string `bun:"end_date,notnull" json:"end_date"`
	Status     enum.SprintStatus `bun:"status,notnull" json:"status"`

	CreateUpdateUnixTimestamp
	SoftDelete
}
