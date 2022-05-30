package example

import (
	"context"
)

type repo struct {
	ctx  *domain.ServiceContext
}

func newRepo(ctx *domain.ServiceContext) *repo {
	return &repo{
		ctx: ctx,
	}
}

func (r *repo) Create(ctx context.Context, example *Example) error {
	return nil
}

func (r *repo) Update(ctx context.Context, example *Example) error {
	return nil
}

func (r *repo) Delete(ctx context.Context, example *Example) error {
	return nil
}

func (r *repo) Get(ctx context.Context, name string) (Example, error) {
	return Example{}, nil
}

func (r *repo) Query(ctx context.Context, param QueryParam) ([]Example, int64, error) {
	result := make([]Example, 0)

	// parse query param
	total := int64(0)

	return result, total, nil
}
