package user

type Service interface {
	GetUserList(page, pageSize int) ([]User, int64, error)
	GetUserDetail(id string) (*User, error)
	AddUser(user *User) error
	DeleteUser(id string) error
	UpdateUser(id string, user *User) (*User, error)
	PartialUpdateUser(id string, updates map[string]interface{}) (*User, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetUserList(page, pageSize int) ([]User, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return s.repo.GetPaged(page, pageSize)
}
func (s *service) GetUserDetail(id string) (*User, error) {
	return s.repo.GetByID(id)
}
func (s *service) AddUser(user *User) error {
	return s.repo.Create(user)
}
func (s *service) DeleteUser(id string) error {
	return s.repo.Delete(id)
}
func (s *service) UpdateUser(id string, user *User) (*User, error) {
	existing, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	user.ID = existing.ID
	if err := s.repo.Update(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *service) PartialUpdateUser(id string, updates map[string]interface{}) (*User, error) {
	if err := s.repo.PartialUpdate(id, updates); err != nil {
		return nil, err
	}
	return s.repo.GetByID(id)
}
