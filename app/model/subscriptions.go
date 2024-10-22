package model

import (
	"app/app/enum"

	"github.com/uptrace/bun"
)

type Subscription struct {
	bun.BaseModel `bun:"table:subscriptions"`

	ID        int64                   `bun:",pk,autoincrement" json:"id"`    // ใช้ ID สำหรับ Primary Key
	UserID    int64                   `bun:"user_id,notnull" json:"user_id"` // FK ใช้ชื่อปกติ
	PlanType  enum.PlanType           `bun:"plan_type,notnull" json:"plan_type"`
	StartDate string                  `bun:"start_date,notnull" json:"start_date"`
	EndDate   string                  `bun:"end_date,notnull" json:"end_date"`
	Status    enum.SubscriptionStatus `bun:"status,notnull" json:"status"`

	CreateUpdateUnixTimestamp
	SoftDelete
}
