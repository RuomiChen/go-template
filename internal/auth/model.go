package auth

type Auth struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     uint   `json:"role"` // 1admin 或 0user
}

type LoginResponse struct {
	Token string `json:"token"`
}
