package model

import (
	"time"

	"app/app/enum"

	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users"`

	ID                 uint                    `bun:",pk,autoincrement" json:"id"`
	Username           string                  `bun:",unique,notnull" json:"username"`
	Firstname          string                  `json:"firstname"`
	Lastname           string                  `json:"lastname"`
	Nickname           string                  `json:"nickname"`
	CitizenID          string                  `bun:",unique,notnull" json:"citizen_id"`
	Birthdate          time.Time               `json:"birthdate"`
	Gender             enum.Gender             `bun:"column:notnull" json:"gender"`
	Nationality        string                  `json:"nationality"`
	RelationshipStatus enum.RelationshipStatus `bun:"column:notnull" json:"relationship_status"`
	Address1           string                  `json:"address_1"`
	Address2           string                  `json:"address_2"`
	MobileNo           string                  `bun:"type:varchar(10)" json:"mobile_no"`
	Email              string                  `bun:",unique,notnull" json:"email"`
	RoleID             uint                    `json:"role_id"`
	Status             enum.Status             `bun:"column:notnull" json:"status"`
	Password           string                  `json:"password"`
	Points             int64                   `json:"points"`
	CreateUpdateUnixTimestamp
	SoftDelete
}
