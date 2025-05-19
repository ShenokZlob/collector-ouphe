package collector

// CheckUser
type CheckUserRequest struct {
	TelegramID int64 `json:"telegram_id"`
}

type CheckUserResponse struct {
	Token   string `json:"token"`
	Success bool   `json:"success"`
}

// Register
type RegisterRequest struct {
	TelegramID int64  `json:"telegram_id"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name,omitempty"`
	Username   string `json:"username,omitempty"`
}

type RegisterResponse struct {
	Token string `json:"token"`
}
