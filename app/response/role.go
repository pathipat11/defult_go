package response

type RoleResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type RoleListResponse struct {
	Roles      []RoleResponse `json:"roles"`
	Pagination Pagination     `json:"pagination"`
}
