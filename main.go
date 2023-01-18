package main

import (
	"flag"
	"github.com/gookit/config/v2"
	"github.com/mkorman9/tiny"
	"github.com/mkorman9/tiny/tinytcp"
	"github.com/rs/zerolog/log"
)

var AppVersion = "development"

func main() {
	configFilePath := flag.String("config", "./config.yml", "path to config.yml file")
	flag.Parse()

	tiny.Init(&tiny.Config{
		ConfigFiles: []string{*configFilePath},
	})

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

	defer func() {
		log.Info().Msgf("Total read: %f MiB", float64(server.Metrics().TotalRead)/1024/1024)
		log.Info().Msgf("Total written: %f MiB", float64(server.Metrics().TotalWritten)/1024/1024)
	}()

	tiny.StartAndBlock(server)
}

func serve(socket *tinytcp.Socket) tinytcp.PacketHandler {
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
