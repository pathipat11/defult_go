package response

type TeamResponse struct {
	ID        string               `bun:"id" json:"id"`
	TeamName  string               `bun:"team_name" json:"team_name"`
	CreatedBy string               `bun:"created_by" json:"created_by"`
	Member    []TeamMemberResponse `bun:"members" json:"members"`
	CreatedAt string               `bun:"created_at" json:"created_at"`
	UpdatedAt string               `bun:"updated_at" json:"updated_at"`
}

type TeamMemberResponse struct {
	UserID      string `bun:"user_id" json:"user_id"`
	DispalyName string `bun:"display_name" json:"display_name"`
}
