package health

type Service interface {
	Check() error
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{
		repository: repository,
	}
}

func (s *service) Check() error {
	return s.repository.Check()
}
