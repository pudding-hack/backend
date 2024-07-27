package auth

type GetCurrentUserResponse struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	Data       User   `json:"data"`
}

type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Address  string `json:"address"`
	Username string `json:"username"`
	Role     Role   `json:"role"`
}

type Role struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
