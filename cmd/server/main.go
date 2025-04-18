package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"sync"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/joho/godotenv"

	"github-insights-dashboard/internal/conf"
)

var (
	Name = "github-insights-dashboard"
	Version string
	flagconf string
	id, _ = os.Hostname()
)

func init() {
	flag.StringVar(&flagconf, "conf", "configs/", "config path, eg: -conf configs/")
}

func newApp(logger log.Logger, hs *http.Server, gs *grpc.Server) *kratos.App {
	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			hs,
			gs,
		),
	)
}

func main() {
	flag.Parse()
	
	_ = godotenv.Load()
	
	logger := log.With(log.NewStdLogger(os.Stdout),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", id,
		"service.name", Name,
		"service.version", Version,
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID(),
	)
	
	bootstrap, err := loadConfig(flagconf, logger)
	if err != nil {
		panic(err)
	}
	
	app, cleanup, err := wireApp(bootstrap, logger)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg := &sync.WaitGroup{}
	
	wg.Add(1)
	go func(c context.Context) {
		defer wg.Done()
		fmt.Println("Starting HTTP and gRPC servers")
		if rErr := app.Run(); rErr != nil {
			panic(rErr)
		}
	}(ctx)
	
	wg.Wait()
}

func loadConfig(configPath string, logger log.Logger) (*conf.Bootstrap, error) {
	
	bootstrap, err := conf.LoadConfig(configPath, logger)
	if err != nil {
		return nil, err
	}

	
	bootstrap, err = conf.LoadSecrets(bootstrap, configPath, logger)
	if err != nil {
		logger.Log(log.LevelWarn, "msg", "Failed to load some secrets, but continuing", "error", err)
	}
	
	
	if bootstrap.Github == nil || bootstrap.Github.Token == "" || bootstrap.Github.Token == "${GITHUB_TOKEN}" {
		
		githubToken := os.Getenv("GITHUB_TOKEN")
		if githubToken != "" {
			if bootstrap.Github == nil {
				bootstrap.Github = &conf.Github{}
			}
			bootstrap.Github.Token = githubToken
			logger.Log(log.LevelInfo, "msg", "Loaded GitHub token from environment")
		} else {
			logger.Log(log.LevelFatal, "msg", "Missing GitHub token in secrets and environment")
			os.Exit(1)
		}
	}

	return bootstrap, nil
} 