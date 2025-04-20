package main

import (
	"context"
	"flag"
	"fmt"
	"luminex-service/internal/conf"
	"os"
	"sync"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
)

var (
	Name     = "luminex-service"
	Version  string
	flagconf string
	id, _    = os.Hostname()
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
	bootstrap, logger := conf.LoadEnvConfig(&flagconf)
	app, err := injectApp(bootstrap, logger)
	if err != nil {
		panic(err)
	}
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
