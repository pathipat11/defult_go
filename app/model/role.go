package model

type Role struct {
	ID          uint   `bun:",pk,autoincrement" json:"id"`
	Name        string `bun:"name" json:"name"`
	Description string `bun:"description" json:"description"`
	CreateUpdateUnixTimestamp
	SoftDelete
}
