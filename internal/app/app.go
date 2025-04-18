package app

import (
	"context"
	"os"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/http"

	"github-insights-dashboard/internal/conf"
	"github-insights-dashboard/internal/biz"
	"github-insights-dashboard/internal/data"
)


type App struct {
	ctx context.Context
	cancel context.CancelFunc
	httpSrv *http.Server
	log     *log.Helper
}


func New(httpSrv *http.Server, logger log.Logger) *App {
	ctx, cancel := context.WithCancel(context.Background())
	return &App{
		ctx:    ctx,
		cancel: cancel,
		httpSrv: httpSrv,
		log:     log.NewHelper(logger),
	}
}


func (a *App) Run() error {
	a.log.Info("starting app")

	if err := a.httpSrv.Start(a.ctx); err != nil {
		a.log.Errorf("failed to start http server: %v", err)
		return err
	}

	return nil
}


func (a *App) Stop() error {
	a.log.Info("stopping app")
	
	if err := a.httpSrv.Stop(a.ctx); err != nil {
		a.log.Errorf("failed to stop http server: %v", err)
		return err
	}

	a.cancel()
	return nil
}


func NewApp(conf *conf.Bootstrap, httpSrv *http.Server, logger log.Logger) *App {
	return New(httpSrv, logger)
}


func NewGithubRepo(c *conf.Bootstrap, logger log.Logger) biz.GithubRepo {
	dataObj, _, _ := data.NewData(c, logger)
	return data.NewGithubRepo(c.Github.Token, logger, dataObj)
}


func NewLogger() log.Logger {
	logger := log.With(log.NewStdLogger(os.Stdout),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"trace_id", tracing.TraceID(),
		"span_id", tracing.SpanID(),
	)
	return logger
} 