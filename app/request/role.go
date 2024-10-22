package request

type CreateRole struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	UserID      []int64 `json:"user_id"`
}

type UpdateRole struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
