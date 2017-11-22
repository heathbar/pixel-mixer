package main

import (
	"fmt"
	"time"

	"github.com/kellydunn/go-opc"
)

func startOpc(server string, mixerOutput chan *Frame) {
	client := opc.NewClient()

	if err := client.Connect("tcp", server); err != nil {
		fmt.Println("Could not connect to destination OPC server", err)
	}

	go func(client *opc.Client, mixerOutput chan *Frame) {
		frameCounter := 0
		timer := time.NewTicker(time.Second)

		for {
			select {
			case f := <-mixerOutput:
				frameCounter++
				if err := client.Send(&f.Message); err != nil {
					fmt.Println("Couldn't send frame", err)
				}
			case <-timer.C:
				//fmt.Printf("%d fps\n", frameCounter)
				frameCounter = 0
			}
		}
	}(client, mixerOutput)
}

func makeFrame() *Frame {
	msg := opc.NewMessage(0)
	msg.SetLength(uint16(pixelCount * 3))
	return &Frame{*msg}
}

func makeFrameOfColor(c Color) *Frame {
	f := makeFrame()

	for i := 0; i < pixelCount; i++ {
		f.Message.SetPixelColor(i, c.R, c.G, c.B)
	}
	return f
}
