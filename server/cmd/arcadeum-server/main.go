package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/horizon-games/arcadeum/server"
	"github.com/horizon-games/arcadeum/server/config"
	"github.com/horizon-games/arcadeum/server/matcher"
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

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Heartbeat("/ping"))
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`¯\_(ツ)_/¯`))
	})

	// Configure websocket route
	r.With(matcher.AddTokenContext).HandleFunc("/ws", server.HandleConnections)

	server.Start()

	// Start the server on localhost
	log.Printf("ARCADEUM Server started :%d; connect at /ws", cfg.ENV.Port)

	if cfg.ENV.TLSEnabled {
		err = http.ListenAndServeTLS(fmt.Sprintf(":%d", cfg.ENV.Port),
			cfg.ENV.TLSCertFile, cfg.ENV.TLSKeyFile, r)
	} else {
		err = http.ListenAndServe(fmt.Sprintf(":%d", cfg.ENV.Port), r)
	}

	if err != nil {
		log.Panic("ListenAndServe:", err)
	}
}
