package server

import (
	v1_m "yuyan/api/marine/v1"
	v1 "yuyan/api/user/v1"
	"yuyan/internal/conf"
	"yuyan/internal/middleware/auth"
	"yuyan/internal/pkg/jwt"
	"yuyan/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, j *conf.Jwt, user *service.UserService, marine *service.MarineServiceService, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			auth.NewAuthMiddleware(jwt.NewJwtRepo(j)).Handler(),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	v1.RegisterUserServiceHTTPServer(srv, user)
	v1_m.RegisterMarineServiceHTTPServer(srv, marine)
	return srv
}
