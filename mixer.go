package main

import (
	"time"
)

func makeMixer(numberOfInputs int, outputEnabler chan bool, inputSelector chan int, output chan *Frame) *Mixer {
	// allocate space for array
	inputs := make([]chan *Frame, numberOfInputs)

	// create the channel instances
	for i := 0; i < numberOfInputs; i++ {
		inputs[i] = make(chan *Frame)
	}
	return &Mixer{inputs: inputs, output: output, outputEnabler: outputEnabler, outputEnabled: false, inputSelector: inputSelector, selectedInput: 0}
}

func (m Mixer) getInput(n int) *chan *Frame {
	return &m.inputs[n]
}

func (m Mixer) loop() {

	const (
		ALPHA = false
		BRAVO = true
	)
	ticker := time.NewTicker(time.Millisecond * 16) // roughly 60 FPS

	blank := make(chan *Frame, 2)
	m.fader.alpha.input = &blank
	m.fader.bravo.input = &blank

	blackFrame := makeFrameOfColor(Color{0, 0, 0})
	m.fader.alpha.cachedFrame = blackFrame
	m.fader.bravo.cachedFrame = blackFrame

	startFadeToInput := func(input *chan *Frame) {
		if m.fader.state == ALPHA {
			m.fader.bravo.input = input

			// handle case where we are fading from RGB to RGB
			// blank out the current so it doesn't receive any new frames as that is undesireable
			if input == m.fader.alpha.input {
				m.fader.alpha.input = &blank
			}
		} else {
			m.fader.alpha.input = input

			// handle case where we are fading from RGB to RGB
			// blank out the current so it doesn't receive any new frames as that is undesireable
			if input == m.fader.bravo.input {
				m.fader.bravo.input = &blank
			}
		}
		m.fader.isFading = true
		m.fader.progress = 0
	}

	for {
		select {
		case e := <-m.outputEnabler:
			if !e {
				blank <- blackFrame
				startFadeToInput(&blank)
			} else if !m.outputEnabled {
				// TODO: define default input, or remember last state
			}
			m.outputEnabled = e
		case newInput := <-m.inputSelector:
			startFadeToInput(&m.inputs[newInput])
		case f := <-*m.fader.alpha.input:
			m.fader.alpha.cachedFrame = f
			m.output <- blendFrames(m.fader.alpha.cachedFrame, m.fader.bravo.cachedFrame, m.fader.progress, m.fader.state)
		case f := <-*(m.fader.bravo.input):
			m.fader.bravo.cachedFrame = f
			m.output <- blendFrames(m.fader.alpha.cachedFrame, m.fader.bravo.cachedFrame, m.fader.progress, m.fader.state)
		case <-ticker.C:
			if m.fader.isFading {
				m.output <- blendFrames(m.fader.alpha.cachedFrame, m.fader.bravo.cachedFrame, m.fader.progress, m.fader.state)
				if m.fader.progress >= 1 {
					m.fader.isFading = false
					m.fader.progress = 0
					m.fader.state = !m.fader.state

					// detach obsolete input
					if m.fader.state == ALPHA {
						m.fader.bravo.input = &blank
					} else {
						m.fader.alpha.input = &blank
					}
				} else {
					m.fader.progress = m.fader.progress + 0.01
				}
			}
		}
	}
}

func blendFrames(frame1 *Frame, frame2 *Frame, progress float64, reverse bool) *Frame {
	if reverse {
		progress = 1 - progress
	}
	if progress <= 0 {
		return frame1
	} else if progress >= 1 {
		return frame2
	}

	newFrame := makeFrame()
	f1 := frame1.Message.ByteArray()
	f2 := frame2.Message.ByteArray()

	for i := 0; i < pixelCount; i++ {
		ii := i*3 + 4
		r := float64(f1[ii]) + progress*float64(int(f2[ii])-int(f1[ii]))
		g := float64(f1[ii+1]) + progress*float64(int(f2[ii+1])-int(f1[ii+1]))
		b := float64(f1[ii+2]) + progress*float64(int(f2[ii+2])-int(f1[ii+2]))

		newFrame.Message.SetPixelColor(i, uint8(r), uint8(g), uint8(b))
	}

	return newFrame
}
