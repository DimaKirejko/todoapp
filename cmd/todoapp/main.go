package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	core_logger "github.com/DimaKirejko/todoapp/internal/core/logger"
	core_postgres_pool "github.com/DimaKirejko/todoapp/internal/core/repository/postgres/pool"
	core_http_middleware "github.com/DimaKirejko/todoapp/internal/core/transport/http/middleware"
	core_http_server "github.com/DimaKirejko/todoapp/internal/core/transport/http/server"
	users_postgres_repository "github.com/DimaKirejko/todoapp/internal/features/users/repository/postgres"
	users_srvice "github.com/DimaKirejko/todoapp/internal/features/users/srvice"
	users_transport_http "github.com/DimaKirejko/todoapp/internal/features/users/transport/http"
	"go.uber.org/zap"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	defer cancel()

	logger, err := core_logger.NewLoger(core_logger.NewConfigMust())
	if err != nil {
		fmt.Println("failed to init application logger:", err)
		os.Exit(1)
	}

	defer func() {
		_ = logger.Sync()
		logger.Close()
	}()

	logger.Debug("initializing postgres connection pool")
	pool, err := core_postgres_pool.NewConnectionPool(ctx, core_postgres_pool.NewConfigMust())

	if err != nil {
		logger.Fatal("Failed to init postgres connection pool", zap.Error(err))
	}
	defer pool.Close()

	logger.Info("Starting todoapp")

	logger.Debug("initializing feature", zap.String("feature", "users"))
	usersRepository := users_postgres_repository.NewUsersRepository(pool)
	usersService := users_srvice.NewUsersService(usersRepository)
	usersTransportHttp := users_transport_http.NewUsersHTTPHandler(usersService)

	logger.Debug("initializing HTTP server")
	httpServer := core_http_server.NewHTTPServer(
		core_http_server.NewConfigMust(),
		logger,
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Panic(),
		core_http_middleware.Trace(),
	)
	apiVersionRouter := core_http_server.NewAPIVersionRoute(core_http_server.ApiVersion1)
	apiVersionRouter.RegisterRoutes(usersTransportHttp.Routes()...)
	httpServer.RegisterAPIRouters(apiVersionRouter) // 6:17 recheck all

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("HTTP server RUN error", zap.Error(err))
	}
}
