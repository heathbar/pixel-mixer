package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

var pixelCount int

func main() {
	// parse cmd line args
	configFile := flag.String("c", "config.json", "Specify a configuration file. Default: config.json")
	flag.Parse()

	// parse config file
	config := loadConfiguration(*configFile)
	pixelCount = config.PixelCount

	// setup communication channels
	mixerOutputEnabler := make(chan bool)
	mixerInputSelector := make(chan int)
	rgbInputColor := make(chan *Color)
	mixerOutput := make(chan *Frame)

	// start the inputs and outputs
	startMqtt(config, mixerOutputEnabler, mixerInputSelector, rgbInputColor)
	startOpc(config.Opc.DestinationServer, mixerOutput)

	// create the mixer
	mixer := makeMixer(len(config.Inputs)+1, mixerOutputEnabler, mixerInputSelector, mixerOutput)

	// wire up all the configured inputs
	startFrameGenerators(config, mixer, rgbInputColor)

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
