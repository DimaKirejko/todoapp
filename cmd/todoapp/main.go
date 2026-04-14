package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	core_logger "github.com/DimaKirejko/todoapp/internal/core/logger"
	core_pgx_pool "github.com/DimaKirejko/todoapp/internal/core/repository/postgres/pool/pgx"
	core_http_middleware "github.com/DimaKirejko/todoapp/internal/core/transport/http/middleware"
	core_http_server "github.com/DimaKirejko/todoapp/internal/core/transport/http/server"
	task_postgres_repository "github.com/DimaKirejko/todoapp/internal/features/tasks/repository/postgres"
	tasks_service "github.com/DimaKirejko/todoapp/internal/features/tasks/service"
	tasks_transport_http "github.com/DimaKirejko/todoapp/internal/features/tasks/transport/http"
	users_postgres_repository "github.com/DimaKirejko/todoapp/internal/features/users/repository/postgres"
	users_srvice "github.com/DimaKirejko/todoapp/internal/features/users/srvice"
	users_transport_http "github.com/DimaKirejko/todoapp/internal/features/users/transport/http"
	"go.uber.org/zap"
)

var (
	timeZone = time.UTC
)

func main() {
	time.Local = timeZone

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

	logger.Debug("application time zone", zap.Any("zone", timeZone))

	logger.Debug("initializing postgres connection pool")
	// pool, err := core_postgres_pool.NewConnectionPool(ctx, core_postgres_pool.NewConfigMust())
	pool, err := core_pgx_pool.NewPool(ctx, core_pgx_pool.NewConfigMust())

	if err != nil {
		logger.Fatal("Failed to init postgres connection pool", zap.Error(err))
	}
	defer pool.Close()

	logger.Info("Starting todoapp")

	logger.Debug("initializing feature", zap.String("feature", "users"))
	usersRepository := users_postgres_repository.NewUsersRepository(pool)
	usersService := users_srvice.NewUsersService(usersRepository)
	usersTransportHttp := users_transport_http.NewUsersHTTPHandler(usersService)

	logger.Debug("initializing feature", zap.String("feature", "tasks"))
	tasksRepository := task_postgres_repository.NewTasksRepository(pool)
	tasksService := tasks_service.NewTasksService(tasksRepository)
	tasksTransportHttp := tasks_transport_http.NewTasksHTTPHandler(tasksService)

	logger.Debug("initializing HTTP server")
	httpServer := core_http_server.NewHTTPServer(
		core_http_server.NewConfigMust(),
		logger,
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Trace(),
		core_http_middleware.Panic(), // 10:49
	)
	apiVersionRouter := core_http_server.NewAPIVersionRoute(core_http_server.ApiVersion1)
	apiVersionRouter.RegisterRoutes(usersTransportHttp.Routes()...)
	apiVersionRouter.RegisterRoutes(tasksTransportHttp.Routes()...)

	apiVersionRouter2 := core_http_server.NewAPIVersionRoute(core_http_server.ApiVersion2, core_http_middleware.Dummy("api v2 middleware"))
	apiVersionRouter2.RegisterRoutes(usersTransportHttp.Routes()...)

	httpServer.RegisterAPIRouters(apiVersionRouter, apiVersionRouter2) // 6:17 recheck all

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("HTTP server RUN error", zap.Error(err))
	} //12:21 перевір як працює з кастомними мідлвеа core_http_middleware.Dummy //12:35
}
