package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

func main() {
	configFile := flag.String("c", "config.json", "Specify a configuration file. Default: config.json")
	flag.Parse()

	config := loadConfiguration(*configFile)

	fmt.Println(config.Mqtt.Server)
	fmt.Println(config.Opc.DestinationServer)

	// don't stop believing
	<-make(chan int)
}

func loadConfiguration(file string) Config {
	var config Config
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	return config
}
