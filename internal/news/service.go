package news

type Service interface {
	GetNewsList() ([]News, error)
	AddNews(news *News) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetNewsList() ([]News, error) {
	return s.repo.GetAll()
}

func (s *service) AddNews(news *News) error {
	return s.repo.Create(news)
}
