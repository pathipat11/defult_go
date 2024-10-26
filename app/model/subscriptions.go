package model

import (
	"app/app/enum"

	"github.com/uptrace/bun"
)

type Subscription struct {
	bun.BaseModel `bun:"table:subscriptions"`

	ID        string                  `bun:",default:gen_random_uuid(),pk" json:"id"`
	UserID    string                  `bun:"user_id,notnull" json:"user_id"` // FK ใช้ชื่อปกติ
	PlanType  enum.PlanType           `bun:"plan_type,notnull" json:"plan_type"`
	StartDate int64                   `bun:"start_date,notnull" json:"start_date"`
	EndDate   int64                   `bun:"end_date,notnull" json:"end_date"`
	Status    enum.SubscriptionStatus `bun:"status,notnull" json:"status"`

	CreateUpdateUnixTimestamp
	SoftDelete
}
