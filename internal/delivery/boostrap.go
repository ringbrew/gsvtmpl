package delivery

import (
	"demo/internal/domain"
	"github.com/ringbrew/gsv/server"
	"github.com/ringbrew/gsv/service"
)

func NewRpcServer(ctx *domain.ServiceContext) server.Server {
	opt := server.Classic()

	// set the server port
	opt.Port = ctx.Config.Port.Rpc
	opt.ProxyPort = ctx.Config.Port.Proxy

	// set the trace info
	opt.TraceOption.Sampler = 1
	opt.TraceOption.Exporter = ctx.Config.Trace.Type
	opt.TraceOption.Endpoint = ctx.Config.Trace.Endpoint

	return server.NewServer(server.GRPC, &opt)
}

func RpcService(ctx *domain.ServiceContext) []service.Service {
	return nil
}
