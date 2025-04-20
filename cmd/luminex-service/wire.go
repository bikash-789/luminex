//go:build wireinject
// +build wireinject

package main

import (
	"github-insights-dashboard/internal/biz"
	"github-insights-dashboard/internal/conf"
	"github-insights-dashboard/internal/data"
	"github-insights-dashboard/internal/server"
	"github-insights-dashboard/internal/service"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

func wireApp(conf *conf.Bootstrap, logger log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(
		server.ProviderSet,
		data.ProviderSet,
		biz.ProviderSet,
		service.ProviderSet,
		provideToken,
		newApp,
	))
}

func provideToken(conf *conf.Bootstrap) string {
	return conf.Github.Token
}
