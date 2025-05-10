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
	exportcontroller "homework/internal/export/controller"
	exportservice "homework/internal/export/service"
	lotteryservice "homework/internal/lottery/service"
	"homework/internal/models"
	resultcontroller "homework/internal/result/controller"
	resultrepository "homework/internal/result/repository"
	resultservice "homework/internal/result/service"
	"homework/internal/server"
	"homework/internal/storage"
	"homework/pkg/errors"
	"homework/pkg/log"
	"os"
	"os/signal"

	paymentcontroller "homework/internal/payment/controller"
	paymentrepository "homework/internal/payment/repository"
	paymentservice "homework/internal/payment/service"

	ticketcontroller "homework/internal/ticket/controller"
	ticketrepository "homework/internal/ticket/repository"
	ticketservice "homework/internal/ticket/service"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func runMigrations(dsn string, log log.Logger) {
	m, err := migrate.New("file://migrations", dsn)
	if err != nil {
		log.Error("migration error", "error", err)
		os.Exit(1)
	}

	err = m.Up()

	switch {
	case errors.Is(err, migrate.ErrNoChange):
		log.Info("migration not needed, schema in actual state")
		return

	case err != nil:
		log.Error("migration failed", "error", err)
		os.Exit(1)
	}

	log.Info("migration success")
}

func main() {
	configPath := flag.String("config", "", "configuration file. Can be YAML or JSON file")
	flag.Parse()

	// Загрузка конфига.
	cfg := start(config.NewConfigFromFile(*configPath))
	cfg = config.EnvEnrichment(cfg)

	logger := start(log.New(cfg.Logger))

	// Менеджер сервисов.
	ctx := context.Background()
	manager := NewManager(&ctx, logger)

	dsn := storage.BuildDSN(&cfg.Storage.Postgres)
	runMigrations(dsn, logger)

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
	lotteryLog := serverlog.WithGroup("lottery")

	lotteryService := start(lotteryservice.NewLotteryService(
		lotteryservice.WithLogger(lotteryLog.WithGroup("service")),
	))

	lotteryService.RegisterLottery(models.NewLottery5from36())
	lotteryService.RegisterLottery(models.NewLottery6from45())

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
		drawservice.WithLotteryService(lotteryService),
	))

	// Инициализация контроллера Draw.
	drawController := start(drawcontroller.NewHandler(
		drawcontroller.WithLogger(drawlog.WithGroup("controller")),
		drawcontroller.WithService(drawService),
	))

	resultLog := serverlog.WithGroup("result")

	// Инициализация репозитория Draw result.
	resultRepo := start(resultrepository.NewRepository(
		resultrepository.WithStorage(st),
		resultrepository.WithLogger(resultLog.WithGroup("repository")),
	))

	// Инициализация сервиса DrawResult.
	resultService := start(resultservice.NewResultService(
		resultservice.WithResultLogger(resultLog.WithGroup("service")),
		resultservice.WithResultRepository(resultRepo),
		resultservice.WithLotteryService(lotteryService),
		resultservice.WithDrawService(drawService),
	))

	// Инициализация контроллера DrawResult.
	resultController := start(resultcontroller.NewHandler(
		resultcontroller.WithLogger(resultLog.WithGroup("controller")),
		resultcontroller.WithService(resultService),
	))

	// Родительский логгер для подсистем внутри сервиса ticket.
	ticketlog := serverlog.WithGroup("ticket")

	// Инициализация репозитория Ticket.
	ticketRepo := start(ticketrepository.NewRepository(
		ticketrepository.WithStorage(st),
		ticketrepository.WithLogger(ticketlog.WithGroup("repository")),
	))

	// Инициализация сервиса Ticket.
	ticketService := start(ticketservice.NewTicketService(
		ticketservice.WithTicketLogger(ticketlog.WithGroup("service")),
		ticketservice.WithTicketRepository(ticketRepo),
		ticketservice.WithLotteryService(lotteryService),
		ticketservice.WithDrawService(drawService),
	))

	ticketService.StartExpiredTicketsCleaner(ctx)

	// Инициализация контроллера Ticket.
	ticketController := start(ticketcontroller.NewHandler(
		ticketcontroller.WithLogger(ticketlog.WithGroup("controller")),
		ticketcontroller.WithService(ticketService),
	))

	// Родительский логгер для подсистем внутри сервиса ticket.
	paymentLog := serverlog.WithGroup("payment")

	// Инициализация репозитория Payment.
	paymentRepo := start(paymentrepository.NewRepository(
		paymentrepository.WithStorage(st),
		paymentrepository.WithLogger(paymentLog.WithGroup("repository")),
	))

	// Инициализация сервиса Payment.
	paymentService := start(paymentservice.NewPaymentService(
		paymentservice.WithPaymentLogger(paymentLog.WithGroup("service")),
		paymentservice.WithPaymentRepository(paymentRepo),
		paymentservice.WithTicketService(ticketService),
		paymentservice.WithDrawService(drawService),
	))

	// Инициализация контроллера Payment.
	paymentController := start(paymentcontroller.NewHandler(
		paymentcontroller.WithLogger(paymentLog.WithGroup("controller")),
		paymentcontroller.WithService(paymentService),
	))

	// Родительский логгер для подсистем внутри сервиса export.
	exportLog := serverlog.WithGroup("export")

	// Инициализация сервиса Export.
	exportService := start(exportservice.NewExportService(
		exportservice.WithExportLogger(exportLog),
		exportservice.WithDrawService(drawService),
		exportservice.WithResultService(resultService),
	))

	// Инициализация контроллера Export.
	exportController := start(exportcontroller.NewHandler(
		exportcontroller.WithLogger(exportLog.WithGroup("controller")),
		exportcontroller.WithService(exportService),
	))

	// Инициализация HTTP сервера.
	http := start(server.New(cfg.Server.HTTP,
		server.WithLogger(serverlog.WithGroup("server")),
		server.WithController(authController),
		server.WithController(drawController),
		server.WithController(resultController),
		server.WithController(ticketController),
		server.WithController(paymentController),
		server.WithController(exportController),
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
