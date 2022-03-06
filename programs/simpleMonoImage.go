package main

import (
	"fmt"
	"os"

	"github.com/mipsmonsta/epd"
	"github.com/mipsmonsta/epd/epd_config"
	"github.com/mipsmonsta/epd/imageutil"
)

func main() {

	e := epd.Epd{
		Config: epd_config.EpdConfig{},
	}
	e.Setup()
	e.Clear()

	img, err := imageutil.OpenImage("../imageutil/test/test.jpg")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	e.Display(&img, epd.MODE_MONO_DITHER_ON)

	e.Sleep()

}
