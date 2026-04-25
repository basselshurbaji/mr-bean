package bean

import "context"

type BeanService interface {
	ListBeans(ctx context.Context, userID string) ([]Bean, error)
	CreateBean(ctx context.Context, userID string, p BeanParams) (*Bean, error)
	UpdateBean(ctx context.Context, id, userID string, p BeanParams) (*Bean, error)
	DeleteBean(ctx context.Context, id, userID string) error
}

type beanService struct {
	repo BeanRepo
}

func NewBeanService(repo BeanRepo) BeanService {
	return &beanService{repo: repo}
}

func (s *beanService) ListBeans(ctx context.Context, userID string) ([]Bean, error) {
	return s.repo.ListBeans(ctx, userID)
}

func (s *beanService) CreateBean(ctx context.Context, userID string, p BeanParams) (*Bean, error) {
	return s.repo.CreateBean(ctx, userID, p)
}

func (s *beanService) UpdateBean(ctx context.Context, id, userID string, p BeanParams) (*Bean, error) {
	return s.repo.UpdateBean(ctx, id, userID, p)
}

func (s *beanService) DeleteBean(ctx context.Context, id, userID string) error {
	return s.repo.DeleteBean(ctx, id, userID)
}
