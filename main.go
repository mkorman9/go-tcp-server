package main

import (
	"flag"
	"fmt"
	"github.com/mkorman9/go-commons/configutil"
	"github.com/mkorman9/go-commons/coreutil"
	"github.com/mkorman9/go-commons/logutil"
	"github.com/mkorman9/go-commons/tcputil"
	"github.com/rs/zerolog/log"
	"os"
)

var AppVersion = "development"

func main() {
	configFilePath := flag.String("config", "./config.yml", "path to config.yml file")
	flag.Parse()

	c, err := configutil.LoadConfigFromFile(*configFilePath)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to load configuration file: %v\n", err)
		os.Exit(1)
	}

	logutil.SetupLogger(
		logutil.Level(c.String("log.level", "info")),
	)

	log.Info().Msgf("Version: %s", AppVersion)

	server := tcputil.NewServer(
		tcputil.Address(c.String("tcp.address", "0.0.0.0:7000")),
	)
	server.ForkingStrategy(tcputil.GoroutinePerConnection(
		tcputil.PacketFramingHandler(
			tcputil.SplitBySeparator([]byte{'\n'}),
			1024,
			8192,
			serve,
		),
	))

	coreutil.StartAndBlock(server)
}

func serve(ctx tcputil.PacketFramingContext) {
	socket := ctx.Socket()

	ctx.OnPacket(func(packet []byte) {
		_, err := socket.Write(packet)
		if err != nil {
			if socket.IsClosed() {
				return
			}

			log.Error().Err(err).Msg("Error while writing to client socket")
		}
	})
}
