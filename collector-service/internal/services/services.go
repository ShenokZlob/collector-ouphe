package services

type Service struct {
	repository Repositorer
}

type Repositorer interface {
}

func NewService(repository Repositorer) *Service {
	return &Service{
		repository: repository,
	}
}
