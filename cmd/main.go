package main

import (
	"errors"
	"os"

	"github.com/JonnyShabli/EffectiveMobile/config"
	pkghttp "github.com/JonnyShabli/EffectiveMobile/pkg/http"
	"github.com/JonnyShabli/EffectiveMobile/pkg/logster"
	"github.com/JonnyShabli/EffectiveMobile/pkg/sig"

	"context"
	"flag"

	"golang.org/x/sync/errgroup"
)

const localConfig = "config/local/config_local.yaml"

func main() {
	var appConfig config.Config
	var configFile string

	// читаем флаги запуска
	flag.StringVar(&configFile, "config", localConfig, "Path to the config file")
	flag.Parse()
	err := config.LoadConfig(configFile, &appConfig)
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

	// создаем технический хэндлер(debug и recoverer)
	techHandler := pkghttp.NewHandler("/", pkghttp.DefaultTechOptions())

	g.Go(func() error {
		return logster.LogIfError(
			logger, pkghttp.RunServer(ctx, appConfig.PrivateAddr, logger, techHandler),
			"Tech server error",
		)
	})

	// ждем завершения
	err = g.Wait()
	if err != nil && !errors.Is(err, sig.ErrSignalReceived) {
		logger.WithError(err).Errorf("Exit reason")
	}
}
