package server

import (
	"github-insights-dashboard/internal/conf"
	"github-insights-dashboard/internal/service"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
)

func NewHTTPServer(c *conf.Bootstrap, s *service.LuminexService, logger log.Logger) *http.Server {
	opts := configureServerOptions(c)
	
	srv := http.NewServer(opts...)
	
	registerAPIRoutes(srv, s, logger)
	
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

func registerAPIRoutes(srv *http.Server, s *service.LuminexService, logger log.Logger) {
	handler := NewLuminexHandler(s, logger)
	
	srv.HandlePrefix("/api", handler)
} 