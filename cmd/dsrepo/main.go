package main

import (
	"flag"
	"fmt"
	"log/slog"

	"github.com/davidjspooner/dshttp/pkg/logevent"
	"github.com/davidjspooner/dsrepo/internal/forest"

	_ "github.com/davidjspooner/dsrepo/internal/impl/binary"
	_ "github.com/davidjspooner/dsrepo/internal/impl/container"
	_ "github.com/davidjspooner/dsrepo/internal/impl/tfregistry"

	_ "github.com/davidjspooner/dsfile/pkg/impl/localfs"
	_ "github.com/davidjspooner/dsfile/pkg/impl/s3fs"
)

func main() {

	configPath := flag.String("config", "config.yaml", "Path to the configuration file")
	flag.Parse()

	loghandler := logevent.NewHandler(&slog.HandlerOptions{})
	log := slog.New(loghandler)

	server, err := forest.NewServer(
		forest.WithLogger(log),
		forest.WithConfigFile(*configPath),
	)
	if err != nil {
		fmt.Println("Failed to create server:", err)
		return
	}
	server.ListenAndServe()
}
