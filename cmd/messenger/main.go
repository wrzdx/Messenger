package main

import (
	"context"
	"fmt"
	auth_bcrypt "messenger/internal/core/auth/bcrypt"
	auth_cookie "messenger/internal/core/auth/cookie"
	auth_jwt "messenger/internal/core/auth/jwt"
	config "messenger/internal/core/config"
	logger "messenger/internal/core/logger"
	"messenger/internal/core/postgres"
	http_middleware "messenger/internal/core/transport/http/middleware"
	http_server "messenger/internal/core/transport/http/server"

	auth_postgres_repository "messenger/internal/features/auth/repository/postgres"
	auth_service "messenger/internal/features/auth/service"
	auth_transport_http "messenger/internal/features/auth/transport/http"
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
	cfg := config.NewConfigMust()
	time.Local = cfg.TimeZone

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
		os.Interrupt,
	)
	defer cancel()

	logger, err := logger.NewLogger(logger.NewConfigMust())
	if err != nil {
		fmt.Println("failed to init application logger:", err)
		os.Exit(1)
	}

	defer logger.Close()

	logger.Debug("application time zone", zap.Any("zone", time.Local))
	logger.Debug("initializing postgres connection pool")
	postgresConfig := postgres.NewConfigMust()
	pool, err := postgres.NewPool(
		ctx,
		postgresConfig,
	)
	if err != nil {
		logger.Fatal("failed to init postgres connection pool", zap.Error(err))
	}

	defer pool.Close()

	logger.Debug("initializing feature", zap.String("feature", "auth"))
	usersRepository := users_postgres_repository.NewUsersRepository(pool, postgresConfig.Timeout)
	sessionsRepository := auth_postgres_repository.NewSessionsRepository(pool, postgresConfig.Timeout)
	hasher := auth_bcrypt.NewBcryptHasher()
	jwtProvider := auth_jwt.NewTokenProvider(auth_jwt.NewConfigMust())
	txManager := postgres.NewTransactionManager(pool)
	authConfig := auth_service.AuthConfig{
		AccessTokenTTL: cfg.AccessTokenTTL,
		SessionTTL:     cfg.SessionTTL,
	}
	authService, err := auth_service.NewAuthService(
		usersRepository,
		sessionsRepository,
		hasher,
		jwtProvider,
		txManager,
		authConfig,
	)

	if err != nil {
		logger.Fatal("failed create auth service", zap.Error(err))
	}

	cookieManager := auth_cookie.NewCookieManager(
		authConfig.SessionTTL,
		cfg.Environment.IsProduction(),
		"/api/v1/auth",
	)
	authTransportHTTP := auth_transport_http.NewAuthHTTPHandler(authService, cookieManager)

	logger.Debug("initializing feature", zap.String("feature", "users"))
	usersService := users_service.NewUsersService(usersRepository, sessionsRepository, txManager)
	usersTransportHTTP := users_transport_http.NewUsersHandler(usersService, cookieManager)

	logger.Debug("initializing HTTP server")
	httpConfig := http_server.NewConfigMust()
	authMW := http_middleware.Auth(jwtProvider)
	router := newHTTPRouter(
		logger,
		httpConfig.AllowedOrigins,
		authTransportHTTP.Router(authMW),
		usersTransportHTTP.Router(authMW),
	)
	httpServer := http_server.NewHTTPServer(
		httpConfig,
		logger,
		router,
	)

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("HTTP server run error", zap.Error(err))
	}
}

func newHTTPRouter(
	log *logger.Logger,
	allowedOrigins []string,
	authRouter chi.Router,
	usersRouter chi.Router,
) chi.Router {
	router := chi.NewRouter()
	router.Use(
		http_middleware.CORS(allowedOrigins),
		http_middleware.RequestID(),
		http_middleware.Logging(log),
		http_middleware.Trace(),
		http_middleware.Recovery(),
	)

	routerV1 := chi.NewRouter()
	routerV1.Mount("/auth", authRouter)
	routerV1.Mount("/users", usersRouter)
	router.Mount("/api/v1", routerV1)

	return router
}
