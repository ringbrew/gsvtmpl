package domain

import (
	"context"
	"github.com/ringbrew/gsv/logger"
	"sync"
	"{{.projectName}}/conf"
)

type ServiceContext struct {
	Name      string
	NodeIp    string
	Config    conf.Config
	Logger    logger.Logger
	Signal    context.Context
	cancel    context.CancelFunc
	WaitGroup sync.WaitGroup
}

func (ctx *ServiceContext) Close() {
	if ctx.cancel != nil {
		ctx.cancel()
	}
	ctx.WaitGroup.Wait()
}

var dsc *ServiceContext

func NewServiceContext(c conf.Config) *ServiceContext {
	if defaultServiceContext == nil {

		daoOption := dao.Option{
			Addr:      c.Mysql.Host,
			Database:  c.Mysql.Database,
			UserName:  c.Mysql.UserName,
			Password:  c.Mysql.Password,
			Service:   name,
			ServiceIp: ip,
		}

		if c.Debug {
			daoOption.LogLevel = logger.Info
		}

		d := dao.NewDataAccess(daoOption)

		cc := cache.New(cache.Client{
			Addr:     c.Redis.Host,
			Password: c.Redis.Password,
			DB:       c.Redis.DB,
		})

		logInstance := logger.NewLogger(logger.Option{
			Service:   name,
			ServiceIp: ip,
		})

		dsc = &ServiceContext{
			Name:   name,
			NodeIp: ip,
			Config: c,
			Logger: logInstance,
		}

		dsc.Signal, dsc.cancel = context.WithCancel(context.Background())
	}

	return defaultServiceContext
}
