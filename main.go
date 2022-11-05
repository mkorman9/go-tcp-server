package main

import (
	"flag"
	"fmt"
	"github.com/mkorman9/go-commons/configutil"
	"github.com/mkorman9/go-commons/coreutil"
	"github.com/mkorman9/go-commons/logutil"
	"github.com/mkorman9/go-commons/tcputil"
	"github.com/rs/zerolog/log"
	"io"
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

	logutil.SetupLogger(c)

	log.Info().Msgf("Version: %s", AppVersion)

	server := tcputil.NewServer(c)
	server.ForkingStrategy(tcputil.GoroutinePerConnection(
		tcputil.FramingHandler(
			tcputil.SplitBySeparator([]byte{'\n'}),
			1024,
			8192,
			serve,
		),
	))

	coreutil.StartAndBlock(server)
}

func serve(p tcputil.PacketReader) {
	socket := p.Socket()

	p.OnPacket(func(packet io.Reader) {
		_, err := io.Copy(socket, packet)
		if err != nil {
			if socket.IsClosed() {
				return
			}

			log.Error().Err(err).Msg("Error while writing to client socket")
		}
	})
}
