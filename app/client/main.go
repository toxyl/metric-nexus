package main

import (
	"fmt"
	"os"
	"strings"

	metrics "github.com/toxyl/metric-nexus"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Printf("Usage:    %s [config file] [action] [key] <value>\n", os.Args[0])
		fmt.Printf("Examples: %s config.yaml   CREATE   demo  'a demo key'\n", os.Args[0])
		fmt.Printf("          %s config.yaml   UPDATE   demo  123\n", os.Args[0])
		fmt.Printf("          %s config.yaml   ADD      demo  123\n", os.Args[0])
		fmt.Printf("          %s config.yaml   SUB      demo  123\n", os.Args[0])
		fmt.Printf("          %s config.yaml   INC      demo\n", os.Args[0])
		fmt.Printf("          %s config.yaml   DEC      demo\n", os.Args[0])
		fmt.Printf("          %s config.yaml   READ     demo\n", os.Args[0])
		return
	}
	conf, err := LoadConfig(os.Args[1])
	if err != nil {
		panic(err)
	}

	action := strings.ToUpper(os.Args[2])
	key := os.Args[3]
	val := ""
	if action == "CREATE" || action == "UPDATE" || action == "ADD" || action == "SUB" {
		if len(os.Args) != 5 {
			fmt.Printf("Action %s need a value!\n", action)
			return
		}
		val = os.Args[4]
	}

	client := metrics.NewClient(conf.Host, conf.Port, conf.APIKey, true)

	switch action {
	case "CREATE":
		client.Create(key, val)
	case "UPDATE":
		client.Update(key, val)
	case "ADD":
		client.Add(key, val)
	case "SUB":
		client.Subtract(key, val)
	case "INC":
		client.Increment(key)
	case "DEC":
		client.Decrement(key)
	case "READ":
		v, err := client.Read(key)
		if err != nil {
			panic(err)
		}
		fmt.Println(v)
	}
}
