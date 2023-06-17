package main

import (
	"fmt"
	"time"

	metrics "github.com/toxyl/metric-nexus"
)

var (
	apiKey = "Hello World"
	t      = time.Now()
)

func startServer() {
	server := metrics.NewServer("127.0.0.1", 3000, "/tmp/state.yaml")
	server.AddAPIKey(apiKey)
	panic(server.Start("", "")) // let MetricNexus create a temporary self-signed certificate
}

func startClient() {
	client := metrics.NewClient("127.0.0.1", 3000, apiKey, true)
	if err := client.Create("uptime", "metric server uptime"); err != nil {
		panic(err)
	}

	for {
		if err := client.Update("uptime", time.Since(t).Seconds()); err != nil {
			fmt.Printf("Failed to set uptime metric: %s\n", err.Error())
		}
		time.Sleep(5 * time.Second)
	}
}

func main() {
	go startServer()
	time.Sleep(5 * time.Second)
	startClient()
}
