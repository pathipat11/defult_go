package model

import (
	"time"

	"app/app/enum"

	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users"`

	ID                 int64                   `bun:"id,pk,autoincrement"`
	Username           string                  `bun:"username,unique,notnull" `
	Firstname          string                  `bun:"firstname"`
	Lastname           string                  `bun:"lastname"`
	Nickname           string                  `bun:"nickname"`
	CitizenID          string                  `bun:"citizen_id,unique,notnull"`
	Birthdate          time.Time               `bun:"birthdate"`
	Gender             enum.Gender             `bun:"gender,column:notnull"`
	Nationality        string                  `bun:"nationality"`
	RelationshipStatus enum.RelationshipStatus `bun:"relationship_status,column:notnull"`
	Address1           string                  `bun:"address_1"`
	Address2           string                  `bun:"address_2"`
	MobileNo           string                  `bun:"mobile_no,type:varchar(10)"`
	Email              string                  `bun:"email,unique,notnull"`
	RoleID             int64                   `bun:"role_id"`
	Status             enum.Status             `bun:"status,column:notnull"`
	Password           string                  `bun:"password"`
	Points             int64                   `bun:"points"`
	CreateUpdateUnixTimestamp
	SoftDelete
}
