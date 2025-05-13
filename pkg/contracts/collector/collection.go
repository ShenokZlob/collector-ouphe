package collector

// GetUserWithCollections
type GetCollectionsRequest struct {
	Token string
}

type GetCollectionsResponse struct {
	TelegramID  string                 `json:"telegram_id"`
	Collections []*UsersCollectionsRef `json:"collections,omitempty"`
}

type UsersCollectionsRef struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// CreateCollection
type CreateCollectionRequest struct {
	Token          string
	CollectionName string
}

type CreateCollectionResponse struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
	Name   string `json:"name"`
}

// RenameCollection
type RenameCollectionRequest struct {
	Token        string
	CollectionID string
}

// DeleteCollection
type DeleteCollectionRequest struct {
	Token        string
	CollectionID string
}
