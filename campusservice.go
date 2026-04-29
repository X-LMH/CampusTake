package main

import (
	"CampusTake/common/logxext"
	"CampusTake/internal/config"
	"CampusTake/internal/handler"
	"CampusTake/internal/svc"
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
	logx.SetWriter(&logxext.BizWriter{Out: os.Stdout})

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	server.Use(logxext.HttpLoggerMiddleware())

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
