package main

func rgbFrameGenerator(rgbSelector chan *Color, output chan *Frame) {
	for {
		output <- makeFrameOfColor(*<-rgbSelector)
	}
}
