package main

import "github.com/mkorman9/go-commons/tcputil"

func serve(client *tcputil.ClientSocket) {
	for {
		packet, err := tcputil.ReadSeparatedPacket(client, []byte{'\n'}, 8192)
		if err != nil {
			if client.IsClosed() {
				break
			}
		}

		_, err = client.Write(packet)
		if err != nil {
			if client.IsClosed() {
				break
			}
		}
	}
}
