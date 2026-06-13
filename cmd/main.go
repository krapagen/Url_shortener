package main

import (
	"Url_shortener/internal/config"
	"Url_shortener/internal/http-server/handlers/url/save"
	"Url_shortener/internal/http-server/middleware/logger"
	"Url_shortener/internal/lib/logger/handlers/slogpretty"
	"Url_shortener/internal/lib/logger/sl"
	"Url_shortener/internal/storage/sqlite"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	LocalEnv = "local"
	DevEnv   = "dev"
	ProdEnv  = "prod"
)

func main() {
	//TODO: init config: cleanenv : done
	cfg := config.MustLoadConfig()

	fmt.Println(cfg)

	//TODO: init logger: slog : done
	log := setupLogger(cfg.Env)
	log.Info("app started", slog.String("env", cfg.Env), slog.String("storage_path", cfg.StoragePath))
	log.Debug("debug messages are enabled")

	//TODO: init storage: sqlite3 : done

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}
	//TODO: init router: chi, "chi render"
	router := chi.NewRouter()

	//middleware
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(logger.NewLogger(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	//handlera
	router.Post("/url/save", save.New(log, storage))
	//TODO: run server
	log.Info("starting server", slog.String("address", cfg.Address))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HttpServer.Timeout,
		WriteTimeout: cfg.HttpServer.Timeout,
		IdleTimeout:  cfg.HttpServer.IdleTimeout,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
		os.Exit(1)
	}
	log.Info("server stopped")

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case LocalEnv:
		log = setupPrettySlog()
	case DevEnv:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case ProdEnv:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
