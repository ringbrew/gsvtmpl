package example

import (
	"context"
	"{{.projectName}}/internal/domain"
)

type Service struct {
	ctx  *domain.ServiceContext
	repo *repo
}

type QueryParam struct {
}

func NewService(ctx *domain.ServiceContext) *Service {
	return &Service{
		ctx:  ctx,
		repo: newRepo(ctx),
	}
}

func (s *Service) Create(ctx context.Context, example *Example) error {
	if err := s.repo.Create(ctx, example); err != nil {
		return err
	}

	return nil
}

func (s *Service) Update(ctx context.Context, example *Example) error {
	if err := s.repo.Update(ctx, example); err != nil {
		return err
	}

	return nil
}

func (s *Service) Delete(ctx context.Context, example *Example) error {
	if err := s.repo.Delete(ctx, example); err != nil {
		return err
	}
	return nil
}

func (s *Service) Get(ctx context.Context, name string) (Example, error) {
	result, err := s.repo.Get(ctx, name)
	if err != nil {
		return Example{}, err
	}

	return result, nil
}

func (s *Service) Query(ctx context.Context, param QueryParam) ([]Example, int64, error) {
	result, total, err := s.repo.Query(ctx, param)
	if err != nil {
		return nil, 0, err
	}

	return result, total, nil
}
