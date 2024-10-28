package model

import (
	"app/app/enum"

	"github.com/uptrace/bun"
)

type SprintTask struct {
	bun.BaseModel `bun:"table:sprint_tasks"`

	ID         string          `bun:",default:gen_random_uuid(),pk" json:"id"`
	SprintID   string          `bun:"sprint_id,notnull" json:"sprint_id"`
	TaskName   string          `bun:"task_name,notnull" json:"task_name"`
	TaskDetail string          `bun:"task_detail,notnull" json:"task_detail"`
	AssignedTo string          `bun:"assigned_to" json:"assigned_to"`
	Status     enum.TaskStatus `bun:"status,notnull" json:"status"`

	CreateUpdateUnixTimestamp
}
