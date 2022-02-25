package epd_config

import (
	"log"

	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/conn/v3/physic"
	"periph.io/x/conn/v3/spi"
	"periph.io/x/conn/v3/spi/spireg"
	"periph.io/x/host/v3"
)

const ( //GPIO pins i.e. GPIO17, GPIO25...
	RST_PIN  = "17"
	DC_PIN   = "25"
	CS_PIN   = "8"
	BUSY_PIN = "24"
)


type EpdConfig struct{

	Port spi.PortCloser
	Conn spi.Conn
	Rst gpio.PinIO
	Dc gpio.PinIO
	Cs gpio.PinIO
	Bs gpio.PinIO
}


func (d *EpdConfig) Setup() error {
	_, err := host.Init()

	if err != nil {
		log.Fatal(err)
	}

	// Use spireg SPI port registry to find the first available SPI bus.
	d.Port, err = spireg.Open("")
	if err != nil {
		log.Fatal(err)
	}

	d.Conn, err = d.Port.Connect(400*physic.KiloHertz, spi.Mode0, 8)
	if err != nil {
		log.Fatal(err)
	}

	d.setupPins()

	return nil
	
}

func (d* EpdConfig) setupPins(){
	d.Rst = gpioreg.ByName("GPIO" + RST_PIN)
	d.Dc = gpioreg.ByName("GPIO" + DC_PIN)
	d.Bs = gpioreg.ByName("GPIO" + BUSY_PIN)
	d.Bs.In(gpio.PullUp, gpio.NoEdge)
	d.Cs = gpioreg.ByName("GPIO" + CS_PIN)

}

func (d* EpdConfig) Digital_writeRST(l gpio.Level){
	d.Rst.Out(l)
}

func (d* EpdConfig) Digital_writeDC(l gpio.Level){
	d.Dc.Out(l)
}

func (d* EpdConfig) Digital_writeCS(l gpio.Level){
	d.Cs.Out(l)
}

func (d* EpdConfig) Digital_readBS() gpio.Level{
	return d.Bs.Read()
}

func (d *EpdConfig) Destroy() {
	d.Port.Close()
}

func (d *EpdConfig) WriteBytes(data []byte){
	if err := d.Conn.Tx(data, nil); err != nil{
		log.Println(err)
	}
}