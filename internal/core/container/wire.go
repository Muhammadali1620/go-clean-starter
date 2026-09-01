//go:build wireinject
// +build wireinject

package container

import (
	"new_project/internal/core/config"

	"github.com/google/wire"
	"github.com/uptrace/bun"
)

func InitializeContainer(
	db *bun.DB,
	cfg *config.Config,
) (*Container, error) {
	wire.Build(
		RepositorySet,
		ServiceSet,
		HandlerSet,
		BotHandlerSet,

		wire.Struct(new(Repositories), "*"),
		wire.Struct(new(Services), "*"),
		wire.Struct(new(Handlers), "*"),
		wire.Struct(new(BotHandlers), "*"),
		wire.Struct(new(WorkerHandlers), "*"),
		wire.Struct(new(Container), "*"),
	)
	return &Container{}, nil
}
