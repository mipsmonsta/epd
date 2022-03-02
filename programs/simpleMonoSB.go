package main

import (
	"epd"
	"epd/epd_config"
	"epd/imageutil"
	"fmt"
	"os"
)

func main() {

	e := epd.Epd{
		Config: epd_config.EpdConfig{},
	}
	e.Setup()
	e.Clear()

	img, err := imageutil.OpenImage("../imageutil/test/test_sb.jpg")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	e.Display(&img, epd.MODE_MONO_DITHER_ON)

	e.Sleep()

}
