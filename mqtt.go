package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func startMqtt(config *Config, mixerOutputEnabler chan bool, mixerInputSelector chan int, rgbInputColor chan *Color) {
	messageHandler := func(client mqtt.Client, msg mqtt.Message) {
		message := string(msg.Payload())
		fmt.Printf("%s: %s\n", msg.Topic(), message)

		switch msg.Topic() {
		case config.Mqtt.Topics.Power:
			if message == "ON" {
				mixerOutputEnabler <- true
			} else {
				mixerOutputEnabler <- false
			}
		case config.Mqtt.Topics.Input:
			// TODO: add inputs to config and map input name to input number
		case config.Mqtt.Topics.Color:
			// c, err := parseColor(message)
			// if err != nil {
			// 	fmt.Printf("Parse Error: %s\n", *err)
			// 	return
			// }
			mixerInputSelector <- rgbInput
			// rgbInputColor <- c
		}
	}

	mqtt.ERROR = log.New(os.Stdout, "", 0)
	opts := mqtt.NewClientOptions().AddBroker(config.Mqtt.Server)
	opts.SetKeepAlive(2 * time.Second)
	opts.SetDefaultPublishHandler(messageHandler)
	opts.SetPingTimeout(1 * time.Second)

	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println("MQTT connection error")
		panic(token.Error())
	}

	if token := c.Subscribe(config.Mqtt.Topics.Power, 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	if token := c.Subscribe(config.Mqtt.Topics.Input, 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	if token := c.Subscribe(config.Mqtt.Topics.Color, 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
}

func parseColor(message string) (*Color, *string) {
	rgbStrings := strings.Split(message, ",")

	if len(rgbStrings) != 3 {
		errorMessage := "There were not 3 comma separated values"
		return nil, &errorMessage
	}

	r, err := strconv.Atoi(rgbStrings[0])
	if err != nil {
		errorMessage := "Segment 1 was not a parseable integer"
		return nil, &errorMessage
	}
	g, err := strconv.Atoi(rgbStrings[1])
	if err != nil {
		errorMessage := "Segment 2 was not a parseable integer"
		return nil, &errorMessage
	}
	b, err := strconv.Atoi(rgbStrings[2])
	if err != nil {
		errorMessage := "Segment 3 was not a parseable integer"
		return nil, &errorMessage
	}

	c := Color{uint8(r), uint8(g), uint8(b)}

	return &c, nil
}
