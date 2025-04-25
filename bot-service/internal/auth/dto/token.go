package dto

type RegisterRequest struct {
	TelegramID       int64  `json:"telegram_id"`
	Username         string `json:"username"`
	TelegramNickname string `json:"telegram_nickname"`
}

type RegisterResponse struct {
	Token string `json:"token"`
}
