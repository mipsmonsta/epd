package fontutil

import (
	"errors"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"log"
	"os"
	"strings"

	_ "image/jpeg"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
)


var (
	
	ErrTooBigForScreen = errors.New("content too big")
	ErrContinueNextScreen = errors.New("content spill to next screen")

)

func LoadStandardFont() (font *truetype.Font, err error) {
	font, err = truetype.Parse(goregular.TTF)
	if err != nil {
		return 
	}
	return
	
}

func PrintCenterWhiteTextBlackImage(fontSize float64, imgWidth int, imgHeight int, text string, invertColors bool, debug bool) (img image.Image, spill_text []string, err error){

	f, err := LoadStandardFont() 
	if err != nil{
		log.Fatalf("%s \n", err)
	}

	opts:=truetype.Options{}
	opts.Size = fontSize
	face := truetype.NewFace(f, &opts)

	//check invert colors flag, if true -> black text on white background
	fg, bg := image.White, image.Black

	if invertColors {
		fg, bg = image.Black, image.White
	}

	rgba := image.NewRGBA(image.Rect(0,0, imgWidth, imgHeight))
	//read https://go.dev/blog/image-draw
	draw.Draw(rgba, rgba.Bounds(), bg, image.Pt(0, 0), draw.Src) //draw the background color in rgba image bounds
	c:= freetype.NewContext()
	c.SetFont(f)
	c.SetFontSize(fontSize)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(fg)

	//check if can fit as single line
	textWidth := font.MeasureString(face, text).Ceil(); //round to nearest int
	textHeight := face.Metrics().Ascent.Ceil() + face.Metrics().Descent.Ceil()

	//reduced allowable width for text block
	allowableWidth := imgWidth - 10; //5 on both size
	allowableHeight := imgHeight - 10;

	if textWidth <= allowableWidth {
		x:= (imgWidth - textWidth) / 2
		y:= (imgHeight - textHeight) / 2
		pt := freetype.Pt(x, y)
		_, _ = c.DrawString(text, pt) 

		if textHeight > allowableHeight {
			rgba = nil
			spill_text = append(spill_text, "")
			err = ErrTooBigForScreen
			return
		}
	} else {
		//determine how many rows are needed
		totalRows := 1
		lineString := ""
		lineWidth := 0
		splitStrings := strings.Fields(text) //empty string, so split by each word
		
		//check if any words itself is longer than imgWidth
		for _, splittedText := range(splitStrings){
			strWidth := font.MeasureString(face, splittedText + " ").Ceil()
			if strWidth > allowableWidth {
				rgba = nil
				spill_text = append(spill_text, "")
				err = ErrTooBigForScreen
				return
			}
		}
		
		
		//layout
		for i, segString := range(splitStrings){
			strWidth := font.MeasureString(face, segString + " ").Ceil()
			if strWidth + lineWidth < allowableWidth {
				//stay on line
				lineString += segString + " "
				lineWidth += strWidth
			} else { //new row
				//let write the previous row
				if 5 + totalRows * textHeight > allowableHeight {
					//recover spill over text
					spill_text = strings.Fields(lineString)
					spill_text = append(spill_text, splitStrings[i:]...)
					err = ErrContinueNextScreen

					if debug {
						_ = printDebug(rgba, "./test/test.jpg")
					}	
					
					return
				}
				pt:= freetype.Pt(5, 5 + totalRows * textHeight)
				_, _= c.DrawString(lineString, pt)
				lineWidth = 0
				lineString = segString + " "
				lineWidth += strWidth
				totalRows += 1
			}
		}
		//last row
		pt := freetype.Pt(5, 5 + totalRows * textHeight)
		_, _= c.DrawString(lineString, pt)
	}

	if debug {
		_ = printDebug(rgba, "./test/test.jpg")
	}
	img = rgba
	spill_text = append(spill_text, "")
	err = nil
	return

}

func printDebug(img image.Image, filePath string) error {
	file, err := os.Create(filePath)
		if err != nil{
			return fmt.Errorf("cannot save file to filepath: %w\n", err)
		}
		err = jpeg.Encode(file, img, nil)
		if err != nil {
			return fmt.Errorf("cannot save as jpeg: %w\n", err)
		}

		return nil
}

