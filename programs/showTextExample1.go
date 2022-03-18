package main

import (
	"github.com/mipsmonsta/epd"
	"github.com/mipsmonsta/epd/epd_config"
	"github.com/mipsmonsta/epd/fontutil"
)

func main() {
	text := "Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry."
	img, err := fontutil.PrintCenterWhiteTextBlackImage(20.0, 264, 176, text, true, false) //black text on white 
	if err != nil {
		panic(err)
	}

	e := epd.Epd{
		Config: epd_config.EpdConfig{},
	}
	e.Setup()
	e.Clear()

	e.Display(&img, epd.MODE_MONO_DITHER_OFF)

	e.Sleep()
}