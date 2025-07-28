package main

import (
	"flag"

	"github.com/greencoda/auth0-api-gateway/internal/module"
)

func main() {
	var configFile string

	flag.StringVar(&configFile, "c", "config.yaml", "path to config file")
	flag.Parse()

	module.NewServerModule(configFile).Run()
}
