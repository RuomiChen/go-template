package user

type Service interface {
	GetUsers() ([]User, error)
	AddUser(user *User) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetUsers() ([]User, error) {
	return s.repo.GetAll()
}

func (s *service) AddUser(user *User) error {
	return s.repo.Create(user)
}
