package request

type CreateTeam struct {
	Name string `json:"name" binding:"required"`
}

type UpdateTeam struct {
	Name string `json:"name"`
}

type CreateTeamMember struct {
	TeamID string `json:"team_id" binding:"required"`
	UserID string `json:"user_id" binding:"required"`
}
