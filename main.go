package main

import (
	"github.com/mkorman9/go-commons/configutil"
	"github.com/mkorman9/go-commons/coreutil"
	"github.com/mkorman9/go-commons/logutil"
	"github.com/mkorman9/go-commons/tcputil"
	"github.com/rs/zerolog/log"
)

func main() {
	c := configutil.LoadFromEnv()
	logutil.SetupLogger(c)

	server := tcputil.NewServer(c)
	server.ForkingStrategy(tcputil.GoroutinePerConnection(serve))

	coreutil.StartAndBlock(server)
}

func serve(socket *tcputil.ClientSocket) {
	log.Info().Msgf("Client connected from: %s", socket.RemoteAddress())

	socket.OnClose(func() {
		log.Info().Msgf("Client disconnected: %s", socket.RemoteAddress())
	})

	for {
		data, err := tcputil.ReadSeparatedPacket(socket, []byte{'\n'}, 1024)
		if err != nil {
			if socket.IsClosed() {
				return
			}

			log.Error().Err(err).Msg("Error reading from socket")
			continue
		}

		log.Info().Msgf("Received data: %s", string(data))

		_, err = socket.Write(data)
		if err != nil {
			if socket.IsClosed() {
				return
			}

			log.Error().Err(err).Msg("Error writing to socket")
			continue
		}
	}
}
