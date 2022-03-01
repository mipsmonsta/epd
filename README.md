# E-Paper (Go Port) - epd for the Raspberry PI

## Introduction

I have the e-paper hat for the Raspberry Pi, specifically the [2.7 inch display](https://www.waveshare.com/wiki/2.7inch_e-Paper_HAT).

As I was tickering with Golang as a system programming, I wondered whether I could control the hat using the Go language. After much research, I found the [periph.io](https://periph.io/) project, which supports communications via the SPI interface used by the hat. Hence, I started writing the port, converting from the python library to a Go library.

## Communication via SPI

The connections to the hat for the 2.7 inch are as follows, with SPI related lines and GPIO pin in **bold** and *italics*:

VCC	    3.3V
GND	    GND
DIN	    SPI **MOSI** Pin:**GPIO10**
CLK	    SPI **SCK**  Pin: **GPIO11**
CS	    SPI **chip select** (Low active) Pin: **GPIO8**
DC	    Data/Command control pin (High for data, and low for command) Pin: **GPIO25**
RST	    External reset pin (Low for reset)  Pin: **GPIO17**
BUSY	Busy state output pin (Low for busy) Pin: **GPIO24**

Notice that MISO is not set since no data is pulled from the display to the Raspberry Pi.

For MOSI communication, a byte is sent to the hat at time, whether as command or data. 

## API

epdconfig.go - data struct to keep port, connection and pins' information.

The Digital_writeRST, Digital_writeDC, Digital_writeCS and Digital_readBS are exported outside the go module and 
represents the low-level way to interact with the pins. These should be called by encapsulating functions in the epd library. 
You should not have to use these functions yourself. Setup of the epdconfig should also be called by the epd library Setup function.

epd.go - data struct where you will initiate and use in your program.
imageutil/imageutil.go - where you can use the functions written to manipulate images.

*Sample usage*
>   e := epd.Epd{
>		Config: epd_config.EpdConfig{},
>	}
>	e.Setup() 
>	e.Clear()
>
>	img, err := imageutil.OpenImage("../imageutil/test/test_portrait.jpg")
>	if err != nil {
>		fmt.Println(err)
>		os.Exit(1)
>	}
>
>	e.Display(&img, epd.MODE_MONO_DITHER_ON)
>
>	time.Sleep(5*time.Second)
>
>	img, err = imageutil.OpenImage("../imageutil/test/test_shiba.jpg")
>	if err != nil {
>		fmt.Println(err)
>		os.Exit(1)
>	}
>	
>	e.Display(&img, epd.MODE_MONO_DITHER_OFF)
>
>	e.Sleep()

## Road map of features:
Implemented:
- Display image in monochrome (1 bit black and white) with / without dithering
- [Flody-Steinberg Dithering](https://en.wikipedia.org/wiki/Floyd%E2%80%93Steinberg_dithering)
- Auto detect orientation (portriat or landscape) and fit into the screen size

Yet to implement:
- 4 Shades Greyscale display of image



