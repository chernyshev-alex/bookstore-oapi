package main

import (
	"flag"
	"log"

	"github.com/chernyshev-alex/go-bookstore-oapi/cmd/srv"
	"github.com/chernyshev-alex/go-bookstore-oapi/internal/env"
)

func main() {
	var confPath string
	flag.StringVar(&confPath, "conf", "app-conf.toml", "config file")
	flag.Parse()

	env, err := env.LoadConfig(confPath)
	if err != nil {
		log.Fatal("failed load config")
	}
	srv.StartServer(env)
}
