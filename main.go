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

	mixerOutputEnabler := make(chan bool)
	mixerInputSelector := make(chan int)
	rgbInputColor := make(chan *Color)
	mixerOutput := make(chan *Frame)

	mixer := makeMixer(config, mixerOutputEnabler, mixerInputSelector, mixerOutput)

	startMqtt(config, mixerOutputEnabler, mixerInputSelector, rgbInputColor)

	mixer.loop()
}

func loadConfiguration(file string) *Config {
	var config Config
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	return &config
}
