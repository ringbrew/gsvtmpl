package domain

import (
	"context"
	"{{$.projectName}}/conf"
	"sync"
)

type UseCaseContext struct {
	Config    conf.Config
	Signal    context.Context
	cancel    context.CancelFunc
	WaitGroup sync.WaitGroup
}

func (ctx *UseCaseContext) Watch() {
	ctx.WaitGroup.Add(1)
}

func (ctx *UseCaseContext) Close() {
	if ctx.cancel != nil {
		ctx.cancel()
	}
	ctx.WaitGroup.Wait()
}

var dsc *UseCaseContext

func NewUseCaseContext(c conf.Config) *UseCaseContext {
	if dsc == nil {
		dsc = &UseCaseContext{
			Config: c,
		}

		dsc.Signal, dsc.cancel = context.WithCancel(context.Background())
	}

	return dsc
}
