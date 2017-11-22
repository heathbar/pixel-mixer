package main

func makeMixer(config *Config, outputEnabler chan bool, inputSelector chan int, output chan *Frame) *Mixer {
	return &Mixer{output: output, outputEnabler: outputEnabler, outputEnabled: false, inputSelector: inputSelector, selectedInput: 0, progress: 0}
}

func (m Mixer) loop() {
	for {
		select {
		case e := <-m.outputEnabler:
			m.outputEnabled = e
		case <-m.inputSelector:
		}
	}
}
