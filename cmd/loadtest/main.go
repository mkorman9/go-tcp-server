package main

import (
	"bytes"
	"crypto/rand"
	"flag"
	"fmt"
	"github.com/mkorman9/tiny/tinytcp"
	"math"
	"sync"
	"sync/atomic"
	"time"
)

type testMetrics struct {
	minLatency time.Duration
	maxLatency time.Duration
	avgLatency time.Duration
	samples    int64
}

func main() {
	host := flag.String("host", "localhost", "host to use")
	port := flag.Int("port", 7000, "port to use")
	clients := flag.Int("clients", 1, "number of clients to spawn")
	sends := flag.Int("sends", 10, "number of time a client should send payload")
	payloadSize := flag.Int("payload", 512, "size of payload")
	flag.Parse()

	var (
		connectionErrors uint32
		writeErrors      uint32
		readErrors       uint32
		malformedErrors  uint32
		successCount     uint32
	)

	var metricsMutex sync.Mutex
	metrics := testMetrics{
		minLatency: math.MaxInt64,
		maxLatency: math.MinInt64,
		avgLatency: 0,
	}

	payload, err := preparePayload(*payloadSize)
	if err != nil {
		fmt.Printf("Failed to generate payload: %v\n", err)
		return
	}

	fmt.Printf("=> Using %s:%d\n", *host, *port)
	fmt.Printf(
		"=> Spawning %d clients, sending %d packets each with payload size %d bytes\n",
		*clients,
		*sends,
		*payloadSize,
	)

	var waitGroup sync.WaitGroup
	waitGroup.Add(*clients)

	startTime := time.Now()

	for i := 0; i < *clients; i++ {
		go func() {
			var latencies []time.Duration

			defer func() {
				metricsMutex.Lock()

				for _, latency := range latencies {
					metrics.samples += 1

					var samples = time.Duration(metrics.samples)
					metrics.avgLatency = metrics.avgLatency*(samples-1)/samples + (latency / samples)

					if latency > metrics.maxLatency {
						metrics.maxLatency = latency
					} else if latency < metrics.minLatency {
						metrics.minLatency = latency
					}
				}

				metricsMutex.Unlock()

				waitGroup.Done()
			}()

			client, err := tinytcp.Dial(fmt.Sprintf("%s:%d", *host, *port))
			if err != nil {
				atomic.AddUint32(&connectionErrors, 1)
				return
			}

			for i := 0; i < *sends; i++ {
				beforeWrite := time.Now()

				var err error
				_, err = client.Write(payload)
				if err != nil {
					atomic.AddUint32(&writeErrors, 1)
					continue
				}

				afterWrite := time.Now()
				writeElapsed := afterWrite.Sub(beforeWrite)

				var receivedPayload []byte
				receivedPayload, err = tinytcp.ReadBytes(client, *payloadSize)
				if err != nil {
					atomic.AddUint32(&readErrors, 1)
					continue
				}

				afterRead := time.Now()
				readElapsed := afterRead.Sub(afterWrite)

				latency := writeElapsed + readElapsed
				latencies = append(latencies, latency)

				if !bytes.Equal(payload[:len(payload)-1], receivedPayload) {
					atomic.AddUint32(&malformedErrors, 1)
					continue
				}

				atomic.AddUint32(&successCount, 1)
			}
		}()
	}

	waitGroup.Wait()

	endTime := time.Now()
	totalTime := endTime.Sub(startTime)

	fmt.Printf("=> FINISHED\n")

	fmt.Printf("=> STATS\n")
	fmt.Printf("  success:\t%d\n", successCount)
	fmt.Printf("  total:\t%v\n", totalTime)
	fmt.Printf("  average:\t%v\n", metrics.avgLatency)
	fmt.Printf("  throughput:\t%v\n", formatThroughput(float64(*clients**payloadSize**sends*2)/totalTime.Seconds()))
	fmt.Printf("  maximum:\t%v\n", metrics.maxLatency)
	fmt.Printf("  minimum:\t%v\n", metrics.minLatency)

	fmt.Printf("=> ERRORS:\n")
	fmt.Printf("  connect:\t%d\n", connectionErrors)
	fmt.Printf("  read: \t%d\n", readErrors)
	fmt.Printf("  write:\t%d\n", writeErrors)
	fmt.Printf("  malformed:\t%d\n", malformedErrors)
}

func preparePayload(size int) ([]byte, error) {
	payload := make([]byte, size+1)

	_, err := rand.Read(payload)
	if err != nil {
		return nil, err
	}

	for i := range payload {
		if payload[i] == '\n' {
			payload[i] = 0
		}
	}

	payload[len(payload)-1] = '\n'

	return payload, nil
}
