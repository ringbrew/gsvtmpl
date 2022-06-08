package delivery

import (
	"{{$.projectName}}/internal/domain"
	"github.com/ringbrew/gsv/server"
	"github.com/ringbrew/gsv/service"
)

func NewServer(ctx *domain.UseCaseContext) server.Server {
	opt := server.Classic()
	// set the server port
	opt.Port = ctx.Config.Port.Rpc
	opt.ProxyPort = ctx.Config.Port.Proxy

	return server.NewServer(server.GRPC, &opt)
}

func ServiceList(ctx *domain.UseCaseContext) []service.Service {
	return nil
}
