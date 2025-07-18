package main

import (
	"os"

	"EffectiveMobile/config"
	pkghttp "EffectiveMobile/pkg/http"
	"EffectiveMobile/pkg/logster"

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
	logger.Infof("error: %v", err)

}
