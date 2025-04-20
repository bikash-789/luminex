package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"luminex-service/internal/biz"
	gh "luminex-service/internal/biz/github"
	"luminex-service/internal/conf"
	svr "luminex-service/internal/server"
	"luminex-service/internal/service"
)

func injectApp(config *conf.Bootstrap, logger log.Logger) (*kratos.App, error) {
	ghConfigs := service.ProvideGithubConfigs(config)
	ghHandler := gh.NewGithubHandler(logger, ghConfigs)
	iLuminexHandler := biz.NewLuminexServiceHandler(logger)
	luminexService := service.NewLuminexService(
		iLuminexHandler,
		ghHandler,
		logger,
	)
	grpcServer := svr.NewGRPCServer(config, luminexService, logger)
	httpServer := svr.NewHTTPServer(config, luminexService, logger)
	app := newApp(logger, httpServer, grpcServer)
	return app, nil
}
