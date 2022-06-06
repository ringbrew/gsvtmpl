package domain

import (
	"context"
	"{{$.projectName}}/conf"
	"sync"
)

type ServiceContext struct {
	Config    conf.Config
	Signal    context.Context
	cancel    context.CancelFunc
	WaitGroup sync.WaitGroup
}

func (ctx *ServiceContext) Watch() {
	ctx.WaitGroup.Add(1)
}

func (ctx *ServiceContext) Close() {
	if ctx.cancel != nil {
		ctx.cancel()
	}
	ctx.WaitGroup.Wait()
}

var dsc *ServiceContext

func NewServiceContext(c conf.Config) *ServiceContext {
	if dsc == nil {
		dsc = &ServiceContext{
			Config: c,
		}

		dsc.Signal, dsc.cancel = context.WithCancel(context.Background())
	}

	return dsc
}
