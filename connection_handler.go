package main

import (
	"github.com/mkorman9/go-commons/tcputil"
	"github.com/rs/zerolog/log"
)

type connectionHandler struct {
}

func newConnectionHandler() *connectionHandler {
	return &connectionHandler{}
}

func (c *connectionHandler) OnAccept(socket *tcputil.ClientSocket) {
	log.Info().Msgf("Client connected from: %s", socket.RemoteAddress())

	socket.OnClose(func() {
		log.Info().Msgf("Client disconnected: %s", socket.RemoteAddress())
	})
}

func (c *connectionHandler) OnPacket(socket *tcputil.ClientSocket, packetData []byte) {
	log.Info().Msgf("Received: %s", string(packetData))

	_, err := socket.Write(packetData)
	if err != nil {
		if socket.IsClosed() {
			return
		}

		log.Error().Err(err).Msg("Error while writing to client socket")
	}
}
