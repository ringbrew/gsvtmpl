package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"{{.projectName}}/conf"
	"{{.projectName}}/internal/delivery"
	"{{.projectName}}/internal/domain"
)

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	config := flag.String("f", "config.yaml", "config file path")
	flag.Parse()

	// 读取配置
	c, err := conf.Load(*config)
	if err != nil {
		log.Fatal(err.Error())
	}

	// 初始化server
	sc := domain.NewServiceContext(c)
	s := delivery.NewRpcServer(sc)
	svcImpl := delivery.RpcService(sc)
	if svcImpl == nil {
		log.Fatal(errors.New("invalid svc implement"))
	}

	// 注册服务实现
	if err := s.Register(svcImpl); err != nil {
		log.Fatal(err.Error())
	}

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-interrupt
		cancel()
	}()

	s.RunWithGateway(ctx)

	sc.Close()
}
