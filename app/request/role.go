package request

type CreateRole struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type UpdateRole struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
