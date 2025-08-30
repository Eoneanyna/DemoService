//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.
package main

import (
	"demoserveice/internal/biz"
	"demoserveice/internal/conf"
	"demoserveice/internal/data"
	"demoserveice/internal/server"
	"demoserveice/internal/service"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// initApp init kratos application.
func initApp(*conf.Server, *conf.Data, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, service.ProviderSet, biz.ProvwireiderSet, data.ProviderSet, newApp))
}
