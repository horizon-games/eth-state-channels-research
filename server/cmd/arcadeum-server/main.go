package main

import (
	"flag"
	"log"
	"os"

	"github.com/horizon-games/arcadeum/server"
	"github.com/horizon-games/arcadeum/server/config"
)

var (
	flags      = flag.NewFlagSet("arcadeum-server", flag.ExitOnError)
	configFile = flags.String("config", "", "path to config file")
)

func main() {
	flags.Parse(os.Args[1:])
	cfg := &config.Config{}
	err := config.NewFromFile(*configFile, os.Getenv("CONFIG"), cfg)
	if err != nil {
		log.Panic(err)
	}

	server, err := server.New(cfg)
	if err != nil {
		log.Panic(err)
	}

	err = server.Start()
	if err != nil {
		log.Panic("ListenAndServe:", err)
	}
}
