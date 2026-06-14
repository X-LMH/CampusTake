package main

import (
	"CampusTake/internal/config"
	"CampusTake/internal/handler"
	"CampusTake/internal/mqs"
	"CampusTake/internal/svc"
	"CampusTake/pkg/logger"
	"flag"
	"fmt"
	"os"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/config.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	logx.MustSetup(c.Log)
	logx.SetWriter(&logger.BizWriter{Out: os.Stdout})

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	if c.Performance.HTTPLog {
		server.Use(logger.HttpLoggerMiddleware())
	}

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	mqs.StartAllConsumers(ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
