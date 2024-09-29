package response

type UserResponse struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Nickname  string `json:"nickname"`
	Email     string `json:"email"`
	RoleID    uint   `json:"role_id"`
	Status    string `json:"status"`
}

type UserListResponse struct {
	Users      []UserResponse `json:"users"`
	Pagination Pagination     `json:"pagination"`
}

type GetUserDetail struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	Firstname string `json:"firstname"`
	Nickname  string `json:"nickname"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email"`
	RoleID    uint   `json:"role_id"`
	Point     int64  `json:"points"`
}
