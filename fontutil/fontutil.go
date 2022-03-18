package fontutil

import (
	"errors"
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


var ErrToBigForScreen = errors.New("content too big")

func LoadStandardFont() (font *truetype.Font, err error) {
	font, err = truetype.Parse(goregular.TTF)
	if err != nil {
		return 
	}
	return
	
}

func PrintCenterWhiteTextBlackImage(fontSize float64, imgWidth int, imgHeight int, text string, invertColors bool, debug bool) (image.Image, error){

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
			return nil, ErrToBigForScreen
		}
	} else {
		//determine how many rows are needed
		totalRows := 1
		lineString := ""
		lineWidth := 0
		splitStrings := strings.Fields(text) //empty string, so split by each word
		
		if len(splitStrings) == 1 { //text is a very long string, cannot fit screen width
			return nil, ErrToBigForScreen
		}
		
		//layout
		for _, segString := range(splitStrings){
			strWidth := font.MeasureString(face, segString + " ").Ceil()
			if strWidth + lineWidth < imgWidth {
				//stay on line
				lineString += segString + " "
				lineWidth += strWidth
			} else { //new row
				//let write the previous row
				if 5 + totalRows * textHeight > allowableHeight {
					return nil, ErrToBigForScreen
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
		pt:= freetype.Pt(5, 5 + totalRows * textHeight)
		_, _= c.DrawString(lineString, pt)
	}

	if debug {
		file, err := os.Create("./test/test.jpg")
		if err != nil{
			log.Fatalf("cannot save file to filepath: %s", err)
		}
		err = jpeg.Encode(file, rgba, nil)
		if err != nil {
			log.Fatalf("Cannot save as jpeg: %s\n", err)
		}
	}

	return rgba, nil
}


