package request

import "app/app/enum"

type LoginUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type CreateUser struct {
	Username           string                  `json:"username"`
	Password           string                  `json:"password"`
	Firstname          string                  `json:"firstname"`
	Lastname           string                  `json:"lastname"`
	Nickname           string                  `json:"nickname"`
	CitizenID          string                  `json:"citizen_id"`
	Birthdate          string                  `json:"birthdate"`
	Gender             enum.Gender             `json:"gender"`
	Nationality        string                  `json:"nationality"`
	RelationshipStatus enum.RelationshipStatus `json:"relationship_status"`
	Address1           string                  `json:"address_1"`
	Address2           string                  `json:"address_2"`
	MobileNo           string                  `json:"mobile_no"`
	Email              string                  `json:"email"`
	Status             enum.Status             `json:"status"`
}

type UpdateUser struct {
	Firstname          string                  `json:"firstname"`
	Lastname           string                  `json:"lastname"`
	Nickname           string                  `json:"nickname"`
	Birthdate          string                  `json:"birthdate"`
	Gender             enum.Gender             `json:"gender"`
	Nationality        string                  `json:"nationality"`
	RelationshipStatus enum.RelationshipStatus `json:"relationship_status"`
	Address1           string                  `json:"address_1"`
	Address2           string                  `json:"address_2"`
	MobileNo           string                  `json:"mobile_no"`
	Email              string                  `json:"email"`
	RoleID             uint                    `json:"role_id"`
}
