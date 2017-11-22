package main

import (
	"time"
)

func startFrameGenerators(config *Config, mixer *Mixer, rgbInputColor chan *Color) {
	// Input 0 is always RGB
	go rgbFrameGenerator(rgbInputColor, mixer.inputs[0])

	for i := range config.Inputs {

		switch config.Inputs[i].Type {
		case "rainbow":
			go rainbowFrameGenerator(mixer.inputs[i+1])
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
