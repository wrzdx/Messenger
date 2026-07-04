package main

import (
	"context"
	"fmt"
	core_auth_bcrypt "messenger/internal/core/auth/bcrypt"
	core_config "messenger/internal/core/config"
	core_logger "messenger/internal/core/logger"
	core_pgx_pool "messenger/internal/core/repository/postgres/pool/pgx"
	core_http_middleware "messenger/internal/core/transport/http/middleware"
	core_http_server "messenger/internal/core/transport/http/server"
	users_postgres_repository "messenger/internal/features/users/repository/postgres"
	users_service "messenger/internal/features/users/service"
	users_transport_http "messenger/internal/features/users/transport/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func main() {
	cfg := core_config.NewConfigMust()
	time.Local = cfg.TimeZone

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
		os.Interrupt,
	)
	defer cancel()

	logger, err := core_logger.NewLogger(core_logger.NewConfigMust())
	if err != nil {
		fmt.Println("failed to init application logger:", err)
		os.Exit(1)
	}

	defer logger.Close()

	logger.Debug("application time zone", zap.Any("zone", time.Local))
	logger.Debug("initializing postgres connection pool")
	pool, err := core_pgx_pool.NewPool(
		ctx,
		core_pgx_pool.NewConfigMust(),
	)
	if err != nil {
		logger.Fatal("failed to init postgres connection pool", zap.Error(err))
	}

	defer pool.Close()

	logger.Debug("initializing feature", zap.String("feature", "users"))
	hasher := core_auth_bcrypt.NewBcryptHasher()
	usersRepository := users_postgres_repository.NewUsersRepository(pool)
	usersService := users_service.NewUsersService(usersRepository, hasher)
	usersTranposrtHTTP := users_transport_http.NewUsersHTTPHandler(usersService)
	logger.Debug("initializing HTTP server")
	httpConfig := core_http_server.NewConfigMust()
	router := chi.NewRouter()
	router.Use(
		core_http_middleware.CORS(httpConfig.AllowedOrigins),
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Trace(),
		core_http_middleware.Panic(),
	)

	routerV1 := chi.NewRouter()
	routerV1.Mount("/users", usersTranposrtHTTP.Router())

	router.Mount("/api/v1", routerV1)
	httpServer := core_http_server.NewHTTPServer(
		httpConfig,
		logger,
		router,
	)

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("HTTP server run error", zap.Error(err))
	}
}
