package main

import (
	"bytes"
	"github.com/mkorman9/go-commons/tcputil"
	"github.com/rs/zerolog/log"
	"sync"
)

type connectionHandler struct {
	readBuffers      map[*tcputil.ClientSocket]*bytes.Buffer
	readBuffersMutex sync.Mutex
}

func newConnectionHandler() *connectionHandler {
	return &connectionHandler{
		readBuffers: map[*tcputil.ClientSocket]*bytes.Buffer{},
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
	c.readBuffers[socket] = bytes.NewBuffer(nil)
}

func (c *connectionHandler) OnRead(socket *tcputil.ClientSocket, data []byte) {
	readBuffer := c.readBuffers[socket]
	_, _ = readBuffer.Write(data)
	b := readBuffer.Bytes()

	for {
		index := bytes.Index(b, []byte{'\n'})
		if index == -1 {
			break
		}

		packetData := b[:index]
		c.onPacket(socket, packetData)

		b = b[index+1:]
		c.readBuffers[socket] = bytes.NewBuffer(b)
	}
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
