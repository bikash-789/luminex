package server

import (
	"luminex-service/internal/conf"
	"luminex-service/internal/service"

	v1 "github.com/bikash-789/comm-protos/luminex/v1"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

func NewGRPCServer(c *conf.Bootstrap, s *service.LuminexService, logger log.Logger) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
		),
	}
	if c.Server.Grpc.Network != "" {
		opts = append(opts, grpc.Network(c.Server.Grpc.Network))
	}
	if c.Server.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(c.Server.Grpc.Addr))
	}
	srv := grpc.NewServer(opts...)
	v1.RegisterLuminexServer(srv, s)
	return srv
}
