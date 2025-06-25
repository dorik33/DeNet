package app

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/dorik33/DeNet/internal/config"
	"github.com/dorik33/DeNet/internal/handlers"
	"github.com/dorik33/DeNet/internal/logger"
	"github.com/dorik33/DeNet/internal/middleware/jwt"
	"github.com/dorik33/DeNet/internal/middleware/log"
	"github.com/dorik33/DeNet/internal/repository/store"
	"github.com/dorik33/DeNet/internal/repository/taskrepo"
	"github.com/dorik33/DeNet/internal/repository/userrepo"
	"github.com/dorik33/DeNet/internal/service/user"
	"github.com/go-chi/chi/v5"
)

type App struct {
	logger   *slog.Logger
	cfg      *config.Config
	router   *chi.Mux
	handlers handlers.Handlers
}

func InitApp() *App {
	cfg := config.LoadConfig(".env")
	logger := logger.InitLogger()
	pool, err := store.NewConnection(cfg)
	if err != nil {
		logger.Error("failed to connect to database", slog.String("error", err.Error()))
		os.Exit(1)
	}

	userRepo := userrepo.NewUserRepository(pool, logger)
	taskRepo := taskrepo.NewTaskRepository(pool, logger)

	service := user.NewUserService(userRepo, taskRepo, logger, cfg)

	handlers := handlers.NewHandlers(service, logger)

	app := App{
		logger:   logger,
		cfg:      cfg,
		router:   chi.NewMux(),
		handlers: handlers,
	}

	return &app
}

func (app *App) Run() {
	app.setupRoutes()

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", app.cfg.ServerCfg.HttpPort),
		Handler:      app.router,
		ReadTimeout:  app.cfg.ServerCfg.HttpReadTimeOut,
		WriteTimeout: app.cfg.ServerCfg.HttpWriteTimeOut,
		IdleTimeout:  app.cfg.ServerCfg.HttpIdleTimeOut,
	}

	app.logger.Info("starting server", slog.String("addr", server.Addr))

	if err := server.ListenAndServe(); err != nil {
		app.logger.Error("start server failed", slog.String("error", err.Error()))
		os.Exit(1)
	}
}

func (app *App) setupRoutes() {
	app.router.Group(func(r chi.Router) {
		r.Use(log.LoggingMiddleware(app.logger))
		r.Post("/register", app.handlers.RegisterHandler())
		r.Post("/login", app.handlers.LoginHandler())
		r.Get("/users/leaderboard", app.handlers.LeaderboardHandler())
	})

	app.router.Group(func(r chi.Router) {
		r.Use(log.LoggingMiddleware(app.logger))
		r.Use(jwt.AuthMiddleware(app.logger, app.cfg))
		r.Post("/users/{id}/referrer", app.handlers.SetReferrerHandler())
		r.Get("/users/{id}/status", app.handlers.StatusHandler())
		r.Post("/users/{id}/tasks/complete", app.handlers.CompleteTaskHandler())
	})
}
