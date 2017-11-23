package main

import opc "github.com/kellydunn/go-opc"

const rgbInput = 0

// Config represents the application configuration
type Config struct {
	PixelCount int `json:"pixel-count,omitempty"`
	Mqtt       struct {
		Server string `json:"server"`
		Topics struct {
			Power string `json:"power"`
			Input string `json:"input"`
			Color string `json:"color"`
		}
	}
	Inputs []ConfigInput
	Opc    struct {
		DestinationServer string `json:"destination-server"`
	}
}

// ConfigInput defines an input from application configuration
type ConfigInput struct {
	Type        string `json:"type"`
	MqttMessage string `json:"mqtt-message"`
	Port        int    `json:"port"`
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
	pixelCount    int
	fader         struct {
		isFading bool
		progress float64
		state    bool
		alpha    struct {
			input       *chan *Frame
			cachedFrame *Frame
		}
		bravo struct {
			input       *chan *Frame
			cachedFrame *Frame
		}
	}
}

// OpcForwarderDevice defines an object that an OPC server can send data to
type OpcForwarderDevice struct {
	channel uint8
	output  chan *Frame
}
