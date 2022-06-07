package main

import (
	"context"
	"{{$.projectName}}/conf"
	"{{$.projectName}}/internal/delivery"
	"{{$.projectName}}/internal/domain"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
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
	ucc := domain.NewUseCaseContext(c)
	s := delivery.NewRpcServer(ucc)
	svcImpl := delivery.ServiceList(ucc)

	// 注册服务实现
	for i := range svcImpl {
		if err := s.Register(svcImpl[i]); err != nil {
			log.Fatal(err.Error())
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-interrupt
		cancel()
	}()

	s.Run(ctx)

	ucc.Close()
}
