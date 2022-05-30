package main

import (
	"context"
	"flag"
	"{{.projectName}}/conf"
	"{{.projectName}}/internal/delivery"
	"{{.projectName}}/internal/domain"
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
	sc := domain.NewServiceContext(c)
	d := delivery.NewWebServer(sc, delivery.TypeApi)

	// 注册Handler
	handlers := delivery.ApiHandler(sc)
	for _, v := range handlers {
		d.Register(v)
	}

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-interrupt
		cancel()
	}()

	// 启动服务
	d.Run(ctx)

	sc.Close()
}
