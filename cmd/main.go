package main

import (
	"EffectiveMobile/config"
	"flag"
)

func main() {
	var appConfig config.Config
	var configFile string

	// читаем флаги запуска
	flag.StringVar(&configFile, "config", "config/local/local_config.yaml", "Path to the config file")
	flag.Parse()
	err := config.LoadConfig(configFile, &appConfig)
	if err != nil {
		panic(err)
	}

}
