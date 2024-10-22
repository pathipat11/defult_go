package model

import (
	"app/app/enum"

	"github.com/uptrace/bun"
)

type Ad struct {
	bun.BaseModel `bun:"table:ads"`

	ID            int64                `bun:",pk,autoincrement" json:"id"`                // ใช้ ID สำหรับ Primary Key
	AdvertiserID  int64                `bun:"advertiser_id,notnull" json:"advertiser_id"` // FK ใช้ชื่อปกติ
	Title         string               `bun:"title,notnull" json:"title"`
	Content       string               `bun:"content" json:"content"`
	DisplayStatus enum.AdDisplayStatus `bun:"display_status,notnull" json:"display_status"`
	ClickCount    int64                `bun:"click_count" json:"click_count"`

	CreateUpdateUnixTimestamp
	SoftDelete
}
