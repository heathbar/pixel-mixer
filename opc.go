package main

import (
	"fmt"
	"time"

	"github.com/kellydunn/go-opc"
)

func startOpc(server string, mixerOutput chan *Frame) {
	client := opc.NewClient()
	frameCounter := 0
	timer := time.NewTicker(time.Second)

	if err := client.Connect("tcp", server); err != nil {
		fmt.Println("Could not connect to destination OPC server", err)
	}

	for {
		select {
		case f := <-mixerOutput:
			//fmt.Print(".")
			frameCounter++
			if err := client.Send(&f.Message); err != nil {
				fmt.Println("Couldn't send frame", err)
			}
		case <-timer.C:
			//fmt.Printf("%d fps\n", frameCounter)
			frameCounter = 0
		}
	}
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
