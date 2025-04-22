package main

import (
	"context"
	"flag"
	"fmt"
	authcontroller "homework/internal/auth/controller"
	authrepository "homework/internal/auth/repository"
	authservice "homework/internal/auth/service"
	"homework/internal/config"
	drawcontroller "homework/internal/draw/controller"
	drawrepository "homework/internal/draw/repository"
	drawservice "homework/internal/draw/service"
	"homework/internal/server"
	"homework/internal/storage"
	"homework/pkg/log"
	"os"
	"os/signal"
)

func main() {
	configPath := flag.String("config", "", "configuration file. Can be YAML or JSON file")
	flag.Parse()

	// Загрузка конфига.
	cfg := start(config.NewConfigFromFile(*configPath))

	logger := start(log.New(cfg.Logger))

	// Менеджер сервисов.
	ctx := context.Background()
	manager := NewManager(&ctx, logger)

	logger.Info("server starting", "build", manager.build)
	defer logger.Info("server stopped")

	storagelog := logger.WithGroup("storage")

	// Инициализация хранилища.
	st := start(storage.NewStorage(
		storage.WithConfig(cfg.Storage.Postgres),
		storage.WithLogger(storagelog),
	))

	// Отдельная группа логгеров для серверов
	serverlog := logger.WithGroup("http")

	// Родительский логгер для подсистем внутри сервиса auth.
	authlog := serverlog.WithGroup("auth")

	// Инициализация репозитория auth.
	authRepo := start(authrepository.NewRepository(
		authrepository.WithStorage(st),
		authrepository.WithLogger(authlog.WithGroup("repository")),
	))

	// Инициализация сервиса auth.
	authService := start(authservice.NewAuthService(
		authservice.WithAuthLogger(authlog.WithGroup("service")),
		authservice.WithAuthRepository(authRepo),
	))

	// Инициализация контроллера auth.
	authController := start(authcontroller.NewHandler(
		authcontroller.WithLogger(authlog.WithGroup("controller")),
		authcontroller.WithService(authService),
	))

	// Родительский логгер для подсистем внутри сервиса draw.
	drawlog := serverlog.WithGroup("draw")

	// Инициализация репозитория Draw.
	drawRepo := start(drawrepository.NewRepository(
		drawrepository.WithStorage(st),
		drawrepository.WithLogger(drawlog.WithGroup("repository")),
	))

	// Инициализация сервиса Draw.
	drawService := start(drawservice.NewDrawService(
		drawservice.WithDrawLogger(drawlog.WithGroup("service")),
		drawservice.WithDrawRepository(drawRepo),
		drawservice.WithAuthService(authService),
	))

	// Инициализация контроллера Draw.
	drawController := start(drawcontroller.NewHandler(
		drawcontroller.WithLogger(drawlog.WithGroup("controller")),
		drawcontroller.WithService(drawService),
	))

	// Инициализация HTTP сервера.
	http := start(server.New(cfg.Server.HTTP,
		server.WithLogger(serverlog.WithGroup("server")),
		server.WithController(authController),
		server.WithController(drawController),
	))

	go manager.run(http.ListenAndServe)

	logger.Info("server started")

	select {
	case <-manager.quit: // Ждем пока все сервисы не остановятся.
	case <-sigint(): // Или сигнал Interrupt.
	}

	logger.Info("server stopping")
}

// startErr завершает работу программы с ошибкой, если err != nil.
func startErr(err error, name string) {
	if err != nil {
		fmt.Fprintf(flag.CommandLine.Output(), "Failed to init %s:\n  %v\n\n", name, err)
		flag.Usage()
		os.Exit(1)
	}
}

// start проверяет ошибку, и если она не nil, то завершает программу.
// Это позволяет проводить инициализацию без однотипного кода.
func start[T any](svc T, err error) T {
	name := fmt.Sprintf("%T", svc)
	startErr(err, name)

	return svc
}

// sigint создаёт сигнал, который принимает события [os.Interrupt].
//
//go:noinline
func sigint() <-chan os.Signal {
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)

	return sigint
}
