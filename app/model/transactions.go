package model

import (
	"app/app/enum"

	"github.com/uptrace/bun"
)

type Transaction struct {
	bun.BaseModel `bun:"table:transactions"`

	ID              int64                  `bun:",pk,autoincrement" json:"id"`                    // ใช้ ID สำหรับ Primary Key
	UserID          int64                  `bun:"user_id,notnull" json:"user_id"`                 // FK ใช้ชื่อปกติ
	SubscriptionID  int64                  `bun:"subscription_id,notnull" json:"subscription_id"` // FK ใช้ชื่อปกติ
	Amount          float64                `bun:"amount,notnull" json:"amount"`
	TransactionDate string                 `bun:"transaction_date,notnull" json:"transaction_date"`
	Status          enum.TransactionStatus `bun:"status,notnull" json:"status"`

	CreateUpdateUnixTimestamp
	SoftDelete
}
