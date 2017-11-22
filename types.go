package main

import opc "github.com/kellydunn/go-opc"

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
