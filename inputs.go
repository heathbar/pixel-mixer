package main

import (
	"math"
	"strconv"
	"time"

	opc "github.com/kellydunn/go-opc"
)

func startFrameGenerators(config *Config, mixer *Mixer, rgbInputColor chan *Color) {
	// Input 0 is always RGB
	go rgbFrameGenerator(rgbInputColor, mixer.inputs[0])

	for i := range config.Inputs {

		switch config.Inputs[i].Type {
		case "channel-walk":
			go channelWalkFrameGenerator(mixer.inputs[i+1])
		case "rainbow":
			go rainbowFrameGenerator(mixer.inputs[i+1])
		case "red-wave":
			go redWaveFrameGenerator(mixer.inputs[i+1])
		case "purple-wave":
			go purpleWaveFrameGenerator(mixer.inputs[i+1])
		case "blue-wave":
			go blueWaveFrameGenerator(mixer.inputs[i+1])
		case "opc":
			go opcFrameForwarder(config.Inputs[i].Port, mixer.inputs[i+1])
		}
	}
}

func rgbFrameGenerator(rgbSelector chan *Color, output chan *Frame) {
	for {
		output <- makeFrameOfColor(*<-rgbSelector)
	}
}

func rainbowFrameGenerator(output chan *Frame) {
	t := 0
	for {
		frame := makeFrame()
		baseHue := float64(t) / 2

		for i := 0; i < pixelCount; i++ {
			hue := (float64(i)*.73 + baseHue) / 360

			color := hsbToRgb(hue, 1.0, 1.0)
			frame.Message.SetPixelColor(i, color.R, color.G, color.B)
		}
		output <- frame

		t++

		if t >= 720 {
			t = 0
		}
		time.Sleep(16 * time.Millisecond) // roughly 60fps
	}
}

func hsbToRgb(h, s, v float64) Color {
	var r, g, b float64
	i := int(h * 6)
	f := h*6 - float64(i)
	p := v * (1 - s)
	q := v * (1 - f*s)
	t := v * (1 - (1-f)*s)

	switch i % 6 {
	case 0:
		r = v
		g = t
		b = p
	case 1:
		r = q
		g = v
		b = p
	case 2:
		r = p
		g = v
		b = t
	case 3:
		r = p
		g = q
		b = v
	case 4:
		r = t
		g = p
		b = v
	case 5:
		r = v
		g = p
		b = q
	}

	// rough color leveling
	sum := r + g + b
	r = r * 255 / sum
	g = g * 255 / sum
	b = b * 255 / sum
	return Color{uint8(r), uint8(g), uint8(b)}
}

func channelWalkFrameGenerator(output chan *Frame) {

	channelWalk := func(color Color, output chan *Frame) {
		frame := makeFrameOfColor(Color{0, 0, 0})
		output <- frame

		for i := 0; i < pixelCount; i++ {
			frame.Message.SetPixelColor(i, color.R, color.G, color.B)
			output <- frame
			time.Sleep(25 * time.Millisecond) // roughly 60fps
		}
	}

	for {
		channelWalk(Color{255, 0, 0}, output)
		channelWalk(Color{0, 255, 0}, output)
		channelWalk(Color{0, 0, 255}, output)
	}
}

func redWaveFrameGenerator(output chan *Frame) {
	const DURATION = 10
	const FPS = 60
	const TOTAL_FRAMES = DURATION * FPS
	const DELTA_PER_FRAME = 2 * math.Pi / TOTAL_FRAMES

	for {
		for frameNumber := 0; frameNumber < TOTAL_FRAMES; frameNumber++ {
			f := makeFrame()
			for pixelNumber := 0; pixelNumber < pixelCount; pixelNumber++ {
				t := math.Sin(float64(frameNumber)*DELTA_PER_FRAME + float64(pixelNumber)*0.1)
				r := float64(160) + float64(72)*t
				g := math.Max(0, float64(48)+float64(128)*t)
				b := 0

				f.Message.SetPixelColor(pixelNumber, uint8(r), uint8(g), uint8(b))
			}
			output <- f
			time.Sleep(16 * time.Millisecond) // roughly 60fps
		}
	}
}

func purpleWaveFrameGenerator(output chan *Frame) {
	const DURATION = 10
	const FPS = 60
	const TOTAL_FRAMES = DURATION * FPS
	const DELTA_PER_FRAME = 2 * math.Pi / TOTAL_FRAMES
	f := makeFrame()

	for {
		for frameNumber := 0; frameNumber < TOTAL_FRAMES; frameNumber++ {

			for pixelNumber := 0; pixelNumber < pixelCount; pixelNumber++ {
				t := math.Sin(float64(frameNumber)*DELTA_PER_FRAME + float64(pixelNumber)*0.1)
				r := float64(150) + float64(60)*t
				g := 0
				b := math.Min(255, float64(255)+float64(30)*t)

				f.Message.SetPixelColor(pixelNumber, uint8(r), uint8(g), uint8(b))
			}
			output <- f
			time.Sleep(16 * time.Millisecond) // roughly 60fps
		}
	}
}

func blueWaveFrameGenerator(output chan *Frame) {
	const DURATION = 10
	const FPS = 60
	const TOTAL_FRAMES = DURATION * FPS
	const DELTA_PER_FRAME = 2 * math.Pi / TOTAL_FRAMES
	f := makeFrame()

	for {
		for frameNumber := 0; frameNumber < TOTAL_FRAMES; frameNumber++ {

			for pixelNumber := 0; pixelNumber < pixelCount; pixelNumber++ {
				t := math.Sin(float64(frameNumber)*DELTA_PER_FRAME + float64(pixelNumber)*0.1)
				r := 0
				g := float64(120) + (float64(10)*t)*(float64(10)*t)*-t
				b := float64(120) + (float64(10)*t)*(float64(10)*t)*t

				f.Message.SetPixelColor(pixelNumber, uint8(r), uint8(g), uint8(b))
			}
			output <- f
			time.Sleep(16 * time.Millisecond) // roughly 60fps
		}
	}
}

// Write sends an OPC message to the given device
func (device OpcForwarderDevice) Write(message *opc.Message) error {
	device.output <- &Frame{*message}
	return nil
}

// Channel retrieves the OPC channel number from the given device
func (device OpcForwarderDevice) Channel() uint8 {
	return device.channel
}

func opcFrameForwarder(port int, output chan *Frame) {
	server := opc.NewServer()
	server.RegisterDevice(OpcForwarderDevice{0, output})
	go server.ListenOnPort("tcp", ":"+strconv.Itoa(port))
	server.Process()
}
