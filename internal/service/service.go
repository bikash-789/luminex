package service

import (
	"github.com/google/wire"
	gh "luminex-service/internal/biz/github"
	"luminex-service/internal/conf"
	"luminex-service/internal/interfaces/entity"
)

var ProviderSet = wire.NewSet(
	NewLuminexService,
	wire.Bind(new(gh.GithubHandler), new(*gh.GithubHandler)),
	ProvideGithubConfigs,
)

func ProvideGithubConfigs(bootstrap *conf.Bootstrap) entity.GithubConfig {
	return conf.GetGithubConfig(bootstrap)
}
