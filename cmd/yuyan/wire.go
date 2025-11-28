//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"yuyan/internal/biz"
	"yuyan/internal/conf"
	"yuyan/internal/data"
	"yuyan/internal/pkg/jwt"
	"yuyan/internal/server"
	"yuyan/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(*conf.Server, *conf.Data, *conf.Jwt, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, jwt.ProviderSet, newApp))
}
