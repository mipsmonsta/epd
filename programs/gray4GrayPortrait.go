package main

import (
	"github.com/mipsmonsta/epd"
	"github.com/mipsmonsta/epd/imageutil"

	"fmt"
	"os"

	"github.com/mipsmonsta/epd/epd_config"
)

func main() {

	e := epd.Epd{
		Config: epd_config.EpdConfig{},
	}
	e.Setup_4Gray()
	e.Clear()

	img, err := imageutil.OpenImage("../imageutil/test/test_portrait.jpg")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	e.Display_4Gray(&img)

	e.Sleep()

}
