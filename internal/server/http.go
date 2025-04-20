package server

import (
	pb "github.com/bikash-789/comm-protos/luminex/v1"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
	"luminex-service/internal/conf"
	"luminex-service/internal/service"
)

func NewHTTPServer(c *conf.Bootstrap, s *service.LuminexService, logger log.Logger) *http.Server {
	opts := configureServerOptions(c)

	srv := http.NewServer(opts...)
	pb.RegisterLuminexHTTPServer(srv, s)

	return srv
}

func configureServerOptions(c *conf.Bootstrap) []http.ServerOption {
	opts := []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
		),
	}

	if c.Server.Http.Network != "" {
		opts = append(opts, http.Network(c.Server.Http.Network))
	}

	if c.Server.Http.Addr != "" {
		opts = append(opts, http.Address(c.Server.Http.Addr))
	}

	return opts
}
