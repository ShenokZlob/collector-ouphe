package auth

// CheckUserRequest — запрос для получения пользователя по telegram_id
// @Description Запрос для проверки существующего пользователя по Telegram ID
// @example { "telegram_id": 123456789 }
type CheckUserRequest struct {
	TelegramID int64 `json:"telegram_id" binding:"required" example:"123456789"`
}

// CheckUserResponse — ответ на проверку пользователя
// @Description Ответ с токеном и флагом успеха
// @example { "token": "eyJhbG...", "success": true }
type CheckUserResponse struct {
	Token   string `json:"token" example:"eyJhbG..."`
	Success bool   `json:"success"`
}

// RegisterRequest — данные для регистрации нового пользователя
// @Description Регистрация пользователя по Telegram ID и данным профиля
// @example { "telegram_id": 123456789, "first_name": "Ivan", "last_name": "Ivanov", "username": "ivan123" }
type RegisterRequest struct {
	TelegramID int64  `json:"telegram_id" binding:"required" example:"123456789"`
	FirstName  string `json:"first_name" binding:"required" example:"Ivan"`
	LastName   string `json:"last_name,omitempty" example:"Ivanov"`
	Username   string `json:"username,omitempty" example:"ivan123"`
}

// RegisterResponse — ответ после регистрации
// @Description Ответ с JWT-токеном
// @example { "token": "eyJhbG..." }
type RegisterResponse struct {
	Token string `json:"token" example:"eyJhbG..."`
}
