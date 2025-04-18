package data

import (
	"github-insights-dashboard/internal/conf"
	"github.com/go-kratos/kratos/v2/log"
	"github-insights-dashboard/internal/biz"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewData,
	NewGithubRepo,
	wire.Bind(new(biz.GithubRepo), new(*githubRepo)),
)

type Data struct {
	log *log.Helper
}

func NewData(c *conf.Bootstrap, logger log.Logger) (*Data, func(), error) {
	logHelper := log.NewHelper(logger)
	
	d := &Data{
		log: logHelper,
	}
	
	return d, func() {
		logHelper.Info("closing the data resources")
	}, nil
}

func NewGithubRepo(token string, logger log.Logger, data *Data) *githubRepo {
	return newGithubRepo(token, logger, data)
} 