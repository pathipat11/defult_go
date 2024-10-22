package model

import "github.com/uptrace/bun"

type Setting struct {
	bun.BaseModel `bun:"table:settings"`

	ID           int64  `bun:",pk,autoincrement" json:"id"`            // ใช้ ID สำหรับ Primary Key
	SettingName  string `bun:"setting_name,notnull" json:"setting_name"`
	SettingValue string `bun:"setting_value,notnull" json:"setting_value"`

	CreateUpdateUnixTimestamp
	SoftDelete
}
