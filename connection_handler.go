package main

import (
	"bytes"
	"github.com/mkorman9/go-commons/tcputil"
	"github.com/rs/zerolog/log"
	"sync"
)

type connectionHandler struct {
	readBuffers      map[*tcputil.ClientSocket][]byte
	readBuffersMutex sync.Mutex
}

func newConnectionHandler() *connectionHandler {
	return &connectionHandler{
		readBuffers: map[*tcputil.ClientSocket][]byte{},
	}
}

func (c *connectionHandler) OnAccept(socket *tcputil.ClientSocket) {
	log.Info().Msgf("Client connected from: %s", socket.RemoteAddress())

	socket.OnClose(func() {
		log.Info().Msgf("Client disconnected: %s", socket.RemoteAddress())

		c.readBuffersMutex.Lock()
		defer c.readBuffersMutex.Unlock()
		delete(c.readBuffers, socket)
	})

	c.readBuffersMutex.Lock()
	defer c.readBuffersMutex.Unlock()
	c.readBuffers[socket] = nil
}

func (c *connectionHandler) OnRead(socket *tcputil.ClientSocket, data []byte) {
	if len(c.readBuffers[socket])+len(data) > 8196 {
		c.readBuffers[socket] = nil
		return
	}

	buffer := bytes.Join([][]byte{
		c.readBuffers[socket],
		data,
	}, nil)

	for {
		packetData, other, ok := bytes.Cut(buffer, []byte{'\n'})
		if !ok {
			break
		}

		c.onPacket(socket, packetData)
		buffer = other
	}

	c.readBuffers[socket] = buffer
}

func (c *connectionHandler) onPacket(socket *tcputil.ClientSocket, packetData []byte) {
	log.Info().Msgf("Received: %s", string(packetData))

	_, err := socket.Write(packetData)
	if err != nil {
		if socket.IsClosed() {
			return
		}

		log.Error().Err(err).Msg("Error while writing to client socket")
	}
}
