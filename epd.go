package epd

import (
	"epd/epd_config"
	"epd/imageutil"
	"fmt"
	"image"
	"image/color"
	"log"
	"time"

	"periph.io/x/conn/v3/gpio"
)

const (
	EPD_WIDTH int = 176
	EPD_HEIGHT int = 264
)

var (
	lut_vcom_dc = []byte{
		0x00, 0x00,
		0x00, 0x1A, 0x1A, 0x00, 0x00, 0x01,
		0x00, 0x0A, 0x0A, 0x00, 0x00, 0x08,
		0x00, 0x0E, 0x01, 0x0E, 0x01, 0x10,
		0x00, 0x0A, 0x0A, 0x00, 0x00, 0x08,
		0x00, 0x04, 0x10, 0x00, 0x00, 0x05,
		0x00, 0x03, 0x0E, 0x00, 0x00, 0x0A,
		0x00, 0x23, 0x00, 0x00, 0x00, 0x01,
	}

	lut_ww = []byte{
		0x90, 0x1A, 0x1A, 0x00, 0x00, 0x01,
		0x40, 0x0A, 0x0A, 0x00, 0x00, 0x08,
		0x84, 0x0E, 0x01, 0x0E, 0x01, 0x10,
		0x80, 0x0A, 0x0A, 0x00, 0x00, 0x08,
		0x00, 0x04, 0x10, 0x00, 0x00, 0x05,
		0x00, 0x03, 0x0E, 0x00, 0x00, 0x0A,
		0x00, 0x23, 0x00, 0x00, 0x00, 0x01,
	}

	// R22H    r
	lut_bw = []byte{
		0xA0, 0x1A, 0x1A, 0x00, 0x00, 0x01,
		0x00, 0x0A, 0x0A, 0x00, 0x00, 0x08,
		0x84, 0x0E, 0x01, 0x0E, 0x01, 0x10,
		0x90, 0x0A, 0x0A, 0x00, 0x00, 0x08,
		0xB0, 0x04, 0x10, 0x00, 0x00, 0x05,
		0xB0, 0x03, 0x0E, 0x00, 0x00, 0x0A,
		0xC0, 0x23, 0x00, 0x00, 0x00, 0x01,
	}

	// R23H    w
	lut_bb = []byte{
		0x90, 0x1A, 0x1A, 0x00, 0x00, 0x01,
		0x40, 0x0A, 0x0A, 0x00, 0x00, 0x08,
		0x84, 0x0E, 0x01, 0x0E, 0x01, 0x10,
		0x80, 0x0A, 0x0A, 0x00, 0x00, 0x08,
		0x00, 0x04, 0x10, 0x00, 0x00, 0x05,
		0x00, 0x03, 0x0E, 0x00, 0x00, 0x0A,
		0x00, 0x23, 0x00, 0x00, 0x00, 0x01,
	}
	// R24H    b
	lut_wb = []byte{
		0x90, 0x1A, 0x1A, 0x00, 0x00, 0x01,
		0x20, 0x0A, 0x0A, 0x00, 0x00, 0x08,
		0x84, 0x0E, 0x01, 0x0E, 0x01, 0x10,
		0x10, 0x0A, 0x0A, 0x00, 0x00, 0x08,
		0x00, 0x04, 0x10, 0x00, 0x00, 0x05,
		0x00, 0x03, 0x0E, 0x00, 0x00, 0x0A,
		0x00, 0x23, 0x00, 0x00, 0x00, 0x01,
	}

	//0-3 gray
	gray_lut_vcom = []byte{
    0x00, 0x00,
    0x00, 0x0A, 0x00, 0x00, 0x00, 0x01,
    0x60, 0x14, 0x14, 0x00, 0x00, 0x01,
    0x00, 0x14, 0x00, 0x00, 0x00, 0x01,
    0x00, 0x13, 0x0A, 0x01, 0x00, 0x01,
    0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
    0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
    0x00, 0x00, 0x00, 0x00, 0x00, 0x00,				
	}	

    //R21
    gray_lut_ww =[]byte{
    0x40, 0x0A, 0x00, 0x00, 0x00, 0x01,
    0x90, 0x14, 0x14, 0x00, 0x00, 0x01,
    0x10, 0x14, 0x0A, 0x00, 0x00, 0x01,
    0xA0, 0x13, 0x01, 0x00, 0x00, 0x01,
    0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
    0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
    0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}
    //R22H	r
    gray_lut_bw = []byte {
    0x40, 0x0A, 0x00, 0x00, 0x00, 0x01,
    0x90, 0x14, 0x14, 0x00, 0x00, 0x01,
    0x00, 0x14, 0x0A, 0x00, 0x00, 0x01,
    0x99, 0x0C, 0x01, 0x03, 0x04, 0x01,
    0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
    0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
    0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}	

    //R23H	w
    gray_lut_wb = []byte{
    0x40, 0x0A, 0x00, 0x00, 0x00, 0x01,
    0x90, 0x14, 0x14, 0x00, 0x00, 0x01,
    0x00, 0x14, 0x0A, 0x00, 0x00, 0x01,
    0x99, 0x0B, 0x04, 0x04, 0x01, 0x01,
    0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
    0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
    0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}
    //R24H	b
    gray_lut_bb = []byte{
    0x80, 0x0A, 0x00, 0x00, 0x00, 0x01,
    0x90, 0x14, 0x14, 0x00, 0x00, 0x01,
    0x20, 0x14, 0x0A, 0x00, 0x00, 0x01,
    0x50, 0x13, 0x01, 0x00, 0x00, 0x01,
    0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
    0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
    0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}
)

type Epd struct {
	Config epd_config.EpdConfig
}

func (e *Epd) Reset() {
	e.Config.Digital_writeRST(gpio.High)
	time.Sleep(200*time.Millisecond)
	e.Config.Digital_writeRST(gpio.Low)
	time.Sleep(5*time.Millisecond)
	e.Config.Digital_writeRST(gpio.High)
	time.Sleep(200*time.Millisecond)
}

func (e *Epd) Send_command(command byte){
	e.Config.Digital_writeDC(gpio.Low)
	e.Config.Digital_writeCS(gpio.Low)
	e.Config.WriteBytes([]byte{command})
	e.Config.Digital_writeCS(gpio.High)
}

func (e *Epd) Send_data(data byte){
	e.Config.Digital_writeDC(gpio.High)
	e.Config.Digital_writeCS(gpio.Low)
	e.Config.WriteBytes([]byte{data})
	e.Config.Digital_writeCS(gpio.High)
}

func (e *Epd) ReadBusy(){
	for e.Config.Digital_readBS() == gpio.Low { //low is idle; 1 is busy
		time.Sleep(200*time.Millisecond)
	}
}

func (e *Epd) Set_lut(){
	e.Send_command(0x20) //vcom
	for count:=0; count < 44; count++{
		e.Send_data(lut_vcom_dc[count])
	}
	e.Send_command(0x21) //ww --
	for count:=0; count < 42; count++{
		e.Send_data(lut_ww[count])
	}
	e.Send_command(0x22) //bw r
	for count:=0; count < 42; count++{
		e.Send_data(lut_bw[count])
	}
	e.Send_command(0x23) //wb w
	for count:=0; count < 42; count++{
		e.Send_data(lut_bb[count])
	}
	e.Send_command(0x24) //bb b
	for count:=0; count < 42; count++{
		e.Send_data(lut_wb[count])
	}
}

func (e *Epd) Gray_SetLut(){
	e.Send_command(0x20)
	for count:=0; count < 44; count++{
		e.Send_data(gray_lut_vcom[count])
	}
      
    e.Send_command(0x21)		
	for count:=0; count < 42; count++{ 	//red not use
		e.Send_data(gray_lut_ww[count])
	}			
	
	e.Send_command(0x22)		
	for count:=0; count < 42; count++{ 	//bw r
		e.Send_data(gray_lut_bw[count])
	}					
 
    e.Send_command(0x23)				//wb w
	for count:=0; count < 42; count++{ 	
		e.Send_data(gray_lut_wb[count])
	}	

	e.Send_command(0x24)				//bb b
	for count:=0; count < 42; count++{ 	
		e.Send_data(gray_lut_bb[count])
	}	
  
	e.Send_command((0x25))				//vcom
	for count:=0; count < 42; count++{ 	
		e.Send_data(gray_lut_ww[count])
	}
}

func (e *Epd) Setup(){
	e.Config.Setup()

	e.Reset()
	
	e.Send_command(0x01) // POWER_SETTING
    e.Send_data(0x03) // VDS_EN, VDG_EN
    e.Send_data(0x00) // VCOM_HV, VGHL_LV[1], VGHL_LV[0]
    e.Send_data(0x2b) // VDH
    e.Send_data(0x2b) // VDL
    e.Send_data(0x09) // VDHR
        
    e.Send_command(0x06) // BOOSTER_SOFT_START
    e.Send_data(0x07)
    e.Send_data(0x07)
    e.Send_data(0x17)
        
    // Power optimization
    e.Send_command(0xF8)
    e.Send_data(0x60)
    e.Send_data(0xA5)
        
    // Power optimization
    e.Send_command(0xF8)
    e.Send_data(0x89)
    e.Send_data(0xA5)
        
	// Power optimization
    e.Send_command(0xF8)
    e.Send_data(0x90)
    e.Send_data(0x00)
        
    // Power optimization
    e.Send_command(0xF8)
    e.Send_data(0x93)
    e.Send_data(0x2A)
        
    // Power optimization
    e.Send_command(0xF8)
    e.Send_data(0xA0)
    e.Send_data(0xA5)
        
    // Power optimization
    e.Send_command(0xF8)
    e.Send_data(0xA1)
    e.Send_data(0x00)
        
    // Power optimization
    e.Send_command(0xF8)
    e.Send_data(0x73)
    e.Send_data(0x41)
        
    e.Send_command(0x16) // PARTIAL_DISPLAY_REFRESH
    e.Send_data(0x00)
    e.Send_command(0x04) // POWER_ON
    e.ReadBusy()

    e.Send_command(0x00) // PANEL_SETTING
    e.Send_data(0xAF) // KW-BF   KWR-AF    BWROTP 0f
        
    e.Send_command(0x30) // PLL_CONTROL
    e.Send_data(0x3A)  // 3A 100HZ   29 150Hz 39 200HZ    31 171HZ
    
    e.Send_command(0X50) // VCOM AND DATA INTERVAL SETTING			
    e.Send_data(0x57)
        
    e.Send_command(0x82) // VCM_DC_SETTING_REGISTER
    e.Send_data(0x12)
    e.Set_lut()
}

func (e *Epd) Setup_4Gray(){
	e.Reset()

	 	e.Send_command(0x01)			//POWER SETTING
        e.Send_data (0x03)
        e.Send_data (0x00)    
        e.Send_data (0x2b)															 
        e.Send_data (0x2b)		

        e.Send_command(0x06)         //booster soft start
        e.Send_data (0x07)		//A
        e.Send_data (0x07)		//B
        e.Send_data (0x17)		//C 

        e.Send_command(0xF8)         //boost??
        e.Send_data (0x60)
        e.Send_data (0xA5)

        e.Send_command(0xF8)         //boost??
        e.Send_data (0x89)
        e.Send_data (0xA5)

        e.Send_command(0xF8)         //boost??
        e.Send_data (0x90)
        e.Send_data (0x00)

        e.Send_command(0xF8)         //boost??
        e.Send_data (0x93)
        e.Send_data (0x2A)

        e.Send_command(0xF8)         //boost??
        e.Send_data (0xa0)
        e.Send_data (0xa5)

        e.Send_command(0xF8)         //boost??
        e.Send_data (0xa1)
        e.Send_data (0x00)

        e.Send_command(0xF8)         //boost??
        e.Send_data (0x73)
        e.Send_data (0x41)

        e.Send_command(0x16)
        e.Send_data(0x00)	

        e.Send_command(0x04)
        e.ReadBusy()

        e.Send_command(0x00)			//panel setting
        e.Send_data(0xbf)		//KW-BF   KWR-AF	BWROTP 0f

        e.Send_command(0x30)			//PLL setting
        e.Send_data (0x90)      	//100hz 

        e.Send_command(0x61)			//resolution setting
        e.Send_data (0x00)		//176
        e.Send_data (0xb0)     	 
        e.Send_data (0x01)		//264
        e.Send_data (0x08)

        e.Send_command(0x82)			//vcom_DC setting
        e.Send_data (0x12)

        e.Send_command(0X50)			///VCOM AND DATA INTERVAL SETTING			
        e.Send_data(0x57)
}

func (e *Epd) Clear(){
	e.Send_command(0x10)
	for i:=0; i < EPD_WIDTH  * EPD_HEIGHT / 8; i++{
		e.Send_data(0xFF)
	}
	e.Send_command(0x13)
	for i:=0; i < EPD_WIDTH  * EPD_HEIGHT / 8; i++{
		e.Send_data(0xFF)
	}	
	e.Send_command(0x12)
	e.ReadBusy()

}

func (e *Epd) Sleep(){
	e.Send_command(0x50)
	e.Send_data(0xF7)
	e.Send_command(0x02)
	e.Send_command(0x07)
	e.Send_data(0xA5)

	time.Sleep(2000*time.Millisecond)
	e.Config.Destroy()
}

func (e *Epd) Display(img *image.Image){
	orientAndfittedImage := imageutil.OrientateAndFitImage(img, EPD_WIDTH, EPD_HEIGHT)
	
	monochromeTensor := ConvertImagetoMonochromeEPDTensor(&orientAndfittedImage)
	monochromeBslices:= GetEPDBuffer(monochromeTensor)
	e.Send_command(0x10)

	for i:=0; i < EPD_HEIGHT * EPD_WIDTH / 8; i++ {
		e.Send_data(0xFF)
	}
	e.Send_command(0x13)
	for i:=0; i < EPD_HEIGHT * EPD_WIDTH / 8; i++ {
		e.Send_data(monochromeBslices[i])
	}
	e.Send_command(0x12)
	e.ReadBusy()
	
}

func (e *Epd) Display_4Gray(img *image.Image){
	orientAndfittedImage := imageutil.OrientateAndFitImage(img, EPD_WIDTH, EPD_HEIGHT)
	
	monochromeTensor := ConvertImagetoMonochromeEPDTensor(&orientAndfittedImage)

	///TODO, convert 8 bit grayscale to 2 bit grayscale
	monochromeBslices:= GetEPDBuffer_4Gray(monochromeTensor)
	e.Send_command(0x10)
	for i:=0; i < EPD_HEIGHT * EPD_WIDTH / 4; i++ {
		var temp3 uint8
		for j:=0; j < 2; j++{
			temp1 := monochromeBslices[i*2+j]
			for k:=0; k < 2; k++{
				temp2 := temp1&0xC0 //0xC0 is 1100 0000
				if temp2 == 0xC0{
					temp3 |= 0x01 //white
				}else if temp2 == 0x00{
					temp3 |= 0x00 //black
				}else if temp2 == 0x08 {
					temp3 |= 0x01 //gray1
				}else {//0x40
					temp3 |= 0x00 //gray2
				}

				temp1 <<=2
				temp2 = temp1&0xC0
				if temp2 == 0xC0 { //white
					temp3 |= 0x01
				}else if temp2 == 0x00{ //black
					temp3 |= 0x00
				}else if temp2 == 0x80{
					temp3 |= 0x01 //gray1
				}else{ //0x04
					temp3 |= 0x00 //gray2
				}
				if j!=1 || k!=1{
					temp3 <<= 1
				}
				temp1 <<= 2
			}
		}
		e.Send_data(temp3)
	}

	e.Send_command(0x13)
	for i:=0; i < EPD_HEIGHT * EPD_WIDTH / 4; i++ {
		var temp3 uint8
		for j:=0; j < 2; j++{
			temp1 := monochromeBslices[i*2+j]
			for k:=0; k < 2; k++{
				temp2 := temp1&0xC0 //0xC0 is 1100 0000
				if temp2 == 0xC0{
					temp3 |= 0x01 //white
				}else if temp2 == 0x00{
					temp3 |= 0x00 //black
				}else if temp2 == 0x08 {
					temp3 |= 0x00 //gray1
				}else {//0x40
					temp3 |= 0x01 //gray2
				}

				temp1 <<=2
				temp2 = temp1&0xC0
				if temp2 == 0xC0 { //white
					temp3 |= 0x01
				}else if temp2 == 0x00{ //black
					temp3 |= 0x00
				}else if temp2 == 0x80{
					temp3 |= 0x00 //gray1
				}else{ //0x40
					temp3 |= 0x01 //gray2
				}
				if j!=1 || k!=1{
					temp3 <<= 1
				}
				temp1 <<= 2
			}
		}
		e.Send_data(temp3)
	}
	e.Gray_SetLut()
	e.Send_command(0x12)
	time.Sleep(200*time.Millisecond)
	e.ReadBusy()
	
}

func ConvertImageto4GrayEPDTensor(img *image.Image)(gray [][]uint8){
	p := imageutil.GetImageTensor(*img)

	//convert to greyscale tensor
	intermediateGreyImg := imageutil.ConvertGreyScale(&p)
	
	p = imageutil.GetImageTensor(intermediateGreyImg)

	for x:=0; x < len(p); x++ {
		col := []uint8{}
		
		for y:=0; y < len(p[0]); y++ {
			pix := p[x][y]
			
			originalClr, ok := color.RGBAModel.Convert(pix).(color.RGBA)
			if !ok {
				log.Fatalln("Color is not RGBA format.")
			} 
			clr := originalClr.A >> 4 //divide by 64
			col = append(col, clr)
		}
		gray = append(gray, col)
	}
	return 
}

func ConvertImagetoMonochromeEPDTensor(img *image.Image)(monochrome [][]uint8){
	p := imageutil.GetImageTensor(*img)

	//convert to greyscale tensor
	intermediateGreyImg := imageutil.ConvertGreyScale(&p)
	threshold := computeOstuThreshold(&intermediateGreyImg)
	fmt.Printf("Computed threshold %d \n", threshold)
	//obtained greyscale tensor 
	p = imageutil.GetImageTensor(intermediateGreyImg)
	dithered := fsDitheringGreyTensorWithThreshold(p, threshold)

	for x:=1; x < len(p)+1; x++ {
		col := []uint8{}
		for y:=1; y < len(p[0])+1; y++ {
			pix := dithered[x][y]
			
			if pix < uint16(threshold) {
				col = append(col, uint8(255))
			} else {
				col = append(col, uint8(0))
			}
		}
		monochrome = append(monochrome, col)
	}
	return 
}

func getHistorgramGreyscaleTensor(grey *image.Image) (hist []int, numPixels int){
	pg := *grey
	size := pg.Bounds().Size()
	his := make([]int, 256)
	for x:=0; x < size.X; x++{
		for y:=0; y < size.Y; y++{
			clr := pg.At(x, y)
			originalColor, ok := color.RGBAModel.Convert(clr).(color.RGBA)
			if ok {
				his[int(originalColor.R)]++
			}

		}
	}
	numPixels = size.X * size.Y
	hist = his
	return
}

func computeOstuThreshold(grey *image.Image) (threshold int) {
	hist, numPixels := getHistorgramGreyscaleTensor(grey)

	sum := float64(0);

	for t:=0; t < 256; t++{
		sum += float64(t * hist[t])
	}
	sumB := float64(0)
	wB := 0
	wF := 0
	varMax := float64(0)
	

	for t:=0; t < 256; t++{
		wB += hist[t];
		if wB == 0{
			continue
		}

		wF = numPixels - wB
		if wF == 0{
			break
		}

		sumB += float64(t * hist[t])

		mB := sumB / float64(wB)
		mF := (sum - sumB) / float64(wF)

		//calculte Between Class Variance
		varBetween := float64(wB) * float64(wF) * (mB - mF) * (mB - mF)
	
		if varBetween > varMax {
			varMax = varBetween
			threshold = t
		}
	}
	return 
}

func GetEPDBuffer(monochrome [][]uint8) []byte{
	imgWidth := len(monochrome)
	imgHeight := len(monochrome[0])

	buffLength := EPD_WIDTH / 8 * EPD_HEIGHT
	buf := make([]byte, buffLength)
	for i:=0; i < buffLength; i++{
		buf[i] = 0xff
	}

	if imgWidth == EPD_WIDTH && imgHeight == EPD_HEIGHT {
		//image is verticals
		for y:=0; y < imgHeight; y++{
			for x:=0; x < imgWidth; x++ {
				if monochrome[x][y] == 0 {
					index := (x + y * EPD_WIDTH) / 8
					buf[index] &= ^(0x80 >> (x % 8)) // x bit will be 0 while rest are 1s whick allow masking
				}
			}
		}
	} else if imgWidth == EPD_HEIGHT && imgHeight == EPD_WIDTH  {
		//image is horizontal
		for y:=0; y < imgHeight; y++ {
			for x:=0; x < imgWidth; x++{
				if monochrome[x][y] == 0 {
					newx := y
					newy := EPD_HEIGHT - newx - 1
					index := (newx + newy*EPD_WIDTH) / 8
					buf[index] &= ^(0x80 >> (newx % 8))

				}
			}
		}
	}

	return buf
}

func findClosetPaletteColorWithThreshold(oldPixelColor uint16, threshold int) (newGrey uint16, quant_error uint16){

	if oldPixelColor > uint16(threshold){
		newGrey = 255 //clipped
	} else {
		newGrey = 0
	}
	quant_error = oldPixelColor - newGrey //will always be a positive uint16 or zero
	return
}

func fsDitheringGreyTensorWithThreshold(pixels [][]color.Color, threshold int) (greypixUint16 [][]uint16){ //use uint16 to prevent overflow of uint8
	y_zeros := []uint16{}
	for a:=0; a < len(pixels[0]) + 2; a++{
		y_zeros = append(y_zeros, 0)
	}
	greypixUint16 = append(greypixUint16, y_zeros) //add left zeros

	for x:=0; x < len(pixels); x++{	
		y_pixels := []uint16{}
		y_pixels = append(y_pixels, 0) //add top zeros
		for y:=0; y < len(pixels[0]); y++{
			clr := pixels[x][y]
			originalClr, ok := color.RGBAModel.Convert(clr).(color.RGBA)
			if !ok {
				break
			}
			y_pixels = append(y_pixels, uint16(originalClr.R))
		}
		y_pixels = append(y_pixels, 0) //add bottom zeros
		greypixUint16 = append(greypixUint16, y_pixels)
	}
	
	y_zeros = []uint16{}
	for a:=0; a < len(pixels[0]) + 2; a++{
		y_zeros = append(y_zeros, 0)
	}
	greypixUint16 = append(greypixUint16, y_zeros) //add right zeros


	for x:=1; x < len(pixels) + 1 ; x++{	
		for y:=1; y < len(pixels[0]) + 1; y++{ 
			//https://www.visgraf.impa.br/Courses/ip00/proj/Dithering1/floyd_steinberg_dithering.html
			//https://en.wikipedia.org/wiki/Floyd%E2%80%93Steinberg_dithering
			newClr, quant_error := findClosetPaletteColorWithThreshold(greypixUint16[x][y], threshold)
			greypixUint16[x][y] = newClr
			greypixUint16[x+1][y] += quant_error * uint16(7) / uint16(16)
			if greypixUint16[x+1][y] > 255{
				greypixUint16[x+1][y] = 255

			}
			greypixUint16[x-1][y+1] += quant_error * uint16(3) / uint16(16)
			if greypixUint16[x-1][y+1] > 255{
				greypixUint16[x-1][y+1] = 255

			}
			greypixUint16[x][y+1] += quant_error * uint16(5) / uint16(16)
			if greypixUint16[x][y+1] > 255{
				greypixUint16[x][y+1] = 255

			}
			greypixUint16[x+1][y+1] += quant_error * uint16(1) / uint16(16)
			if greypixUint16[x+1][y+1] > 255{
				greypixUint16[x+1][y+1] = 255

			}
		}
	}

	return // bigger width + 2 and height + 2
}


func GetEPDBuffer_4Gray(imgTensor [][]uint8) []byte{ //one unint8 represents one pixel
	imgWidth := len(imgTensor)
	imgHeight := len(imgTensor[0]) 
	buf := make([]byte, EPD_WIDTH/4 * EPD_HEIGHT)
	for b:=0; b < EPD_WIDTH/4; b++{
		buf[b] = 0xff
	}


	var i int
	if imgWidth == EPD_WIDTH && imgHeight == EPD_HEIGHT{
		//portrait
		for y:=0; y < imgHeight; y++ {
			for x:= 0; x < imgWidth; x++ {
				if imgTensor[x][y] == 0xC0{
					imgTensor[x][y] = 0x80
				} else if imgTensor[x][y] == 0x80{
					imgTensor[x][y] = 0x40
				}
				i += 1
				if (i % 4 == 0){ //every 4 bytes in tensor
					buf[int(x+(y * EPD_HEIGHT)/4)] = (imgTensor[x-3][y]&0xc0) | ((imgTensor[x-2][y]&0xc0) >> 2) | ((imgTensor[x-1][y]&0xc0) >> 4) | ((imgTensor[x][y]&0xc0) >> 6)
				}
			}
		}
	} else if imgWidth == EPD_HEIGHT && imgHeight == EPD_WIDTH {
		//landscape
		for x:=0; x < imgWidth; x++ {
			for y:=0; y < imgHeight; y++{
				newx := y //buff space
				newy := EPD_HEIGHT - x - 1

				if imgTensor[x][y] == 0xC0{
					imgTensor[x][y] = 0x80
				} else if imgTensor[x][y] == 0x80{
					imgTensor[x][y] = 0x40
				}
				i += 1
				if (i % 4 == 0){ //every 4 bytes in tensor
					buf[int(newx + (newy * EPD_HEIGHT)/4)] = (imgTensor[x][y-3]&0xc0) | ((imgTensor[x][y-2]&0xc0) >> 2) | ((imgTensor[x][y-1]&0xc0) >> 4) | ((imgTensor[x][y]&0xc0) >> 6)
				}

			}
		}
	}
	return buf

}


