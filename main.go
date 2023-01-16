package main

import (
	"flag"
	"fmt"
	"github.com/gookit/config/v2"
	"github.com/mkorman9/tiny"
	"github.com/mkorman9/tiny/tinylog"
	"github.com/mkorman9/tiny/tinytcp"
	"github.com/rs/zerolog/log"
	"os"
)

var AppVersion = "development"

func main() {
	configFilePath := flag.String("config", "./config.yml", "path to config.yml file")
	flag.Parse()

	if loaded := tiny.LoadConfig(*configFilePath); !loaded {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to load configuration file\n")
		os.Exit(1)
	}

	tinylog.SetupLogger(
		tinylog.Level(config.String("log.level", "info")),
	)

	log.Info().Msgf("Version: %s", AppVersion)

	server := tinytcp.NewServer(
		config.String("tcp.address", "0.0.0.0:7000"),
	)
	server.ForkingStrategy(tinytcp.GoroutinePerConnection(
		tinytcp.PacketFramingHandler(
			tinytcp.SplitBySeparator([]byte{'\n'}),
			serve,
		),
	))

	tiny.StartAndBlock(server)
}

func serve(socket *tinytcp.ConnectedSocket) tinytcp.PacketHandler {
	// client connected

	return func(packet []byte) {
		_, err := socket.Write(packet)
		if err != nil {
			if socket.IsClosed() {
				return
			}

			log.Error().Err(err).Msg("Error while writing to client socket")
		}
	}
}
