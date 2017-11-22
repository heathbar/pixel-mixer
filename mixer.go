package main

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

	blank := make(chan *Frame)
	m.fader.alpha.input = &blank
	m.fader.bravo.input = &blank

	blackFrame := makeFrameOfColor(Color{0, 0, 0})
	m.fader.alpha.cachedFrame = blackFrame
	m.fader.bravo.cachedFrame = blackFrame

	startFadeToInput := func(input *chan *Frame) {
		if m.fader.state == ALPHA {
			m.fader.bravo.input = input

			// handle case where we are fading to the same input
			// blank out the current so it doesn't receive any new frames as that is undesireable
			if input == m.fader.alpha.input {
				m.fader.alpha.input = &blank
			}
		} else {
			m.fader.alpha.input = input

			// handle case where we are fading to the same input
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
			m.outputEnabled = e
		case newInput := <-m.inputSelector:
			startFadeToInput(&m.inputs[newInput])
		case f := <-*m.fader.alpha.input:
			m.fader.alpha.cachedFrame = f
			m.output <- m.fader.alpha.cachedFrame
		case f := <-*(m.fader.bravo.input):
			m.fader.bravo.cachedFrame = f
			m.output <- m.fader.bravo.cachedFrame
		}
	}
}
