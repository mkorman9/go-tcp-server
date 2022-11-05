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
	"time"
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
	server.ForkingStrategy(tcputil.WorkerPool(
		tcputil.PacketFraming(8196, tcputil.SplitBySeparator([]byte{'\n'}), newConnectionHandler()),
		8,
		1024,
		10*time.Millisecond,
	))

	coreutil.StartAndBlock(server)
}
