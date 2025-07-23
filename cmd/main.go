package main

import (
	"errors"
	"os"

	"github.com/JonnyShabli/EffectiveMobile/config"
	"github.com/JonnyShabli/EffectiveMobile/internal/controller"
	"github.com/JonnyShabli/EffectiveMobile/internal/repository"
	"github.com/JonnyShabli/EffectiveMobile/internal/service"
	pkghttp "github.com/JonnyShabli/EffectiveMobile/pkg/http"
	"github.com/JonnyShabli/EffectiveMobile/pkg/logster"
	"github.com/JonnyShabli/EffectiveMobile/pkg/postgres"
	"github.com/JonnyShabli/EffectiveMobile/pkg/sig"
	"github.com/joho/godotenv"

	_ "github.com/JonnyShabli/EffectiveMobile/docs"

	"context"
	"flag"

	"golang.org/x/sync/errgroup"
)

const localConfig = "config/local/config_local.yaml"

func main() {
	var appConfig config.Config
	var configFile string

	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	// читаем флаги запуска
	flag.StringVar(&configFile, "config", localConfig, "Path to the config file")
	flag.Parse()
	err = config.LoadConfig(configFile, &appConfig)
	if err != nil {
		panic(err)
	}

	// Создаем логер
	logger := logster.New(os.Stdout, appConfig.Log)
	defer func() { _ = logger.Sync() }()

	// создаем errgroup
	g, ctx := errgroup.WithContext(context.Background())

	// Gracefully shutdown
	g.Go(func() error {
		return sig.ListenSignal(ctx, logger)
	})

	logger.Infof("service starting with config %+v", appConfig)

	// собираем зависимости
	dbConn := postgres.NewConn(ctx, logger, appConfig.DB)
	repo := repository.NewStorage(dbConn)
	subService := service.NewSubsService(repo)
	subController := controller.NewSubsHandler(subService, logger)

	// создаем хэндлер
	handler := pkghttp.NewHandler("/", pkghttp.WithLoger(logger), pkghttp.DefaultTechOptions(), controller.WithApiHandler(subController))
	logger.Infof("create and configure handler %+v", handler)

	if appConfig.DB.Migrate {
		g.Go(func() error {
			return logster.LogIfError(
				logger, postgres.MigrateDB(ctx, dbConn.DB, logger), "migration error")
		})
	} else {
		logger.Infof("migration disabled by config: migrate = %+v", appConfig.DB.Migrate)
	}

	g.Go(func() error {
		return logster.LogIfError(
			logger, pkghttp.RunServer(ctx, appConfig.PrivateAddr, logger, handler),
			"Tech server error",
		)
	})

	// ждем завершения
	err = g.Wait()
	if err != nil && !errors.Is(err, sig.ErrSignalReceived) {
		logger.WithError(err).Errorf("Exit reason")
	}
}
