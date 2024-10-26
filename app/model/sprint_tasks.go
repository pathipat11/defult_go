package model

import (
	"app/app/enum"

	"github.com/uptrace/bun"
)

type SprintTask struct {
	bun.BaseModel `bun:"table:sprint_tasks"`

	ID         string          `bun:",default:gen_random_uuid(),pk" json:"id"`
	SprintID   string          `bun:"sprint_id,notnull" json:"sprint_id"` // FK ใช้ชื่อปกติ
	TaskName   string          `bun:"task_name,notnull" json:"task_name"`
	AssignedTo int64           `bun:"assigned_to,notnull" json:"assigned_to"` // FK ใช้ชื่อปกติ
	Status     enum.TaskStatus `bun:"status,notnull" json:"status"`

	CreateUpdateUnixTimestamp
	SoftDelete
}
