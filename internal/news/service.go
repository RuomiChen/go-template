package news

type Service interface {
	GetNewsList(page, pageSize int) ([]News, int64, error)
	AddNews(news *News) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetNewsList(page, pageSize int) ([]News, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return s.repo.GetPaged(page, pageSize)
}

func (s *service) AddNews(news *News) error {
	return s.repo.Create(news)
}
func (s *service) DeleteNews(id uint) error {
	return s.repo.Delete(id)
}
func (s *service) UpdateNews(news *News) error {
	return s.repo.Update(news)
}
