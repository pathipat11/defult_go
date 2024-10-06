package model

type Role struct {
	ID          int64  `bun:",pk,autoincrement" json:"id"`
	Name        string `bun:"name" json:"name"`
	Description string `bun:"description" json:"description"`
	CreateUpdateUnixTimestamp
	SoftDelete
}
