package collections

// CreateCollectionRequest — запрос для создания новой коллекции
// @Description Запрос для создания коллекции с указанным именем
// @example { "name": "My cool collection" }
type CreateCollectionRequest struct {
	Name string `json:"name" binding:"required" example:"My cool collection"`
}

// RenameCollectionRequest — запрос для переименования коллекции
// @Description Запрос для изменения названия коллекции
// @example { "name": "Renamed collection" }
type RenameCollectionRequest struct {
	Name string `json:"name" binding:"required" example:"Renamed collection"`
}

// Collection — модель коллекции в ответах
// @Description Модель коллекции с ID и именем
// @example { "id": "64a9b66b2db8b91234a6e8e3", "name": "My cool collection" }
type Collection struct {
	ID   string `json:"id" example:"64a9b66b2db8b91234a6e8e3"`
	Name string `json:"name" example:"My cool collection"`
}

// ErrorResponse — стандартная структура ошибки
// @Description Структура ответа при ошибке
// @example { "message": "unauthorized" }
type ErrorResponse struct {
	Message string `json:"message" example:"unauthorized"`
}
