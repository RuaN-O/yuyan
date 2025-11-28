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
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(c *conf.Server, j *conf.Jwt, user *service.UserService, marine *service.MarineServiceService, logger log.Logger) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
			auth.NewAuthMiddleware(jwt.NewJwtRepo(j)).Handler(),
		),
	}
	if c.Grpc.Network != "" {
		opts = append(opts, grpc.Network(c.Grpc.Network))
	}
	if c.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(c.Grpc.Addr))
	}
	if c.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.Grpc.Timeout.AsDuration()))
	}
	srv := grpc.NewServer(opts...)
	v1.RegisterUserServiceServer(srv, user)
	//v1.RegisterGreeterServer(srv, greeter)
	v1_m.RegisterMarineServiceServer(srv, marine)
	return srv
}
