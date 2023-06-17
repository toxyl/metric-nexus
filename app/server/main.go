package main

import (
	"fmt"
	"os"

	metrics "github.com/toxyl/metric-nexus"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage:   %s [config file]\n", os.Args[0])
		fmt.Printf("Example: %s config.yaml\n", os.Args[0])
		return
	}
	conf, err := LoadConfig(os.Args[1])
	if err != nil {
		panic(err)
	}

	server := metrics.NewServer(conf.Host, conf.Port, conf.StateFile)
	for _, k := range conf.APIKeys {
		server.AddAPIKey(k)
	}
	panic(server.Start(conf.KeyFile, conf.CertFile))
}
