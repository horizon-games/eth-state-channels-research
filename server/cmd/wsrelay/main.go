package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/horizon-games/arcadeum/server/config"
	"github.com/horizon-games/arcadeum/server/services/matcher"
	"github.com/horizon-games/arcadeum/server/services/wsrelay"
	serviceConfig "github.com/horizon-games/arcadeum/server/services/wsrelay/config"
)

var (
	flags      = flag.NewFlagSet("wsrelay", flag.ExitOnError)
	configFile = flags.String("config", "", "path to config file")
)

func main() {
	flags.Parse(os.Args[1:])
	cfg := &serviceConfig.Config{}
	err := config.NewFromFile(*configFile, os.Getenv("CONFIG"), cfg)
	if err != nil {
		log.Fatal(err)
	}

	if cfg.ENV.WorkingDir == "" {
		// default to GOPATH
		gopath := os.Getenv("GOPATH")
		if gopath == "" {
			gopath = fmt.Sprintf("%s/go", os.Getenv("HOME")) // make best guess
		}
		cfg.ENV.WorkingDir = fmt.Sprintf("%s/%s", gopath, "src/github.com/horizon-games/arcadeum/server")
	}
	server := wsrelay.NewServer(cfg)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Heartbeat("/ping"))
	r.Use(middleware.Recoverer)

	// Create a simple file server
	// TODO: see go-chi/chi/_examples/fileserver
	// or if you want the files available in prod, look at davatar project
	// fs := http.FileServer(http.Dir("./public/wsrelay"))
	// http.Handle("/", fs)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`¯\_(ツ)_/¯`))
	})

	// Configure websocket route
	r.With(matcher.AddTokenContext).HandleFunc("/ws", server.HandleConnections)

	server.Start()

	// Start the server on localhost
	log.Printf("DGAME WSRELAY Server started :%d; connect at /ws", cfg.ENV.Port)

	if cfg.ENV.TLSEnabled {
		err = http.ListenAndServeTLS(fmt.Sprintf(":%d", cfg.ENV.Port),
			cfg.ENV.TLSCertFile, cfg.ENV.TLSKeyFile, r)
	} else {
		err = http.ListenAndServe(fmt.Sprintf(":%d", cfg.ENV.Port), r)
	}

	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
