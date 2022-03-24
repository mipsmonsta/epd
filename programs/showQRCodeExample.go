package main

import (
	"github.com/mipsmonsta/epd"
	"github.com/mipsmonsta/epd/epd_config"
	"github.com/mipsmonsta/epd/imageutil"
)

func main() {
	img, err:= imageutil.PrintQRCodeWithWhiteBgImageWithURL("https://www.arstechnica.com", 264, 176, imageutil.QRMiddle, 5)
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