package main

import opc "github.com/kellydunn/go-opc"

const rgbInput = 0

// Config represents the application configuration
type Config struct {
	Mqtt struct {
		Server string `json:"server"`
		Topics struct {
			Power string `json:"power"`
			Input string `json:"input"`
			Color string `json:"color"`
		}
	}
	Opc struct {
		DestinationServer string `json:"destination-server"`
	}
}

// Color holds RGB pixel values
type Color struct {
	R, G, B uint8
}

// Frame represents a single frame to be sent as an OPC message
type Frame struct {
	Message opc.Message
}

// Mixer is the heart of the app. It holds all of the inputs and mixes them to the output
type Mixer struct {
	inputs        []chan *Frame
	output        chan *Frame
	outputEnabler chan bool
	outputEnabled bool
	inputSelector chan int
	selectedInput int
	progress      float64
}
