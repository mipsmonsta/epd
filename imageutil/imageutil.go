package imageutil

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"os"
	"sync"

	"github.com/disintegration/imaging"
	"golang.org/x/image/draw"
)

//Scale quality
const (
	ScaleBestQ int = iota	
	ScaleBetterQ
	ScaleGoodQ
	ScaleNormalQ 
)


func OpenImage(path string)(image.Image, error){
	f, err := os.Open(path)
	if err != nil{
		fmt.Println(err)
		return nil, err
	}
	defer f.Close()
	
	img, format, err := image.Decode(f)
	if err!= nil {
		e:= fmt.Errorf("error in decoding: %w", err)
		return nil, e
	}
	
	if format != "jpeg" {
		e:=fmt.Errorf("error in image format - not jpeg")
		return nil, e
	}

	return img, nil
}

func GetImageTensor(img image.Image) (pixels [][]color.Color){
	size := img.Bounds().Size()
	for i:=0; i < size.X; i++ {
		var y []color.Color
		for j:=0; j < size.Y; j++ {
			y = append(y, img.At(i,j))
		}
		pixels = append(pixels, y) // 2 by 2 slices where each contains a color.color
	}
	return
}

func GetBackImage(pixels *[][]color.Color) image.Image{
	rect := image.Rect(0, 0, len(*pixels), len((*pixels)[0]))
	newImage := image.NewRGBA(rect)

	for x:=0; x < len(*pixels); x++ {
		for y:=0; y < len((*pixels)[0]); y++ {
			p:=(*pixels)[x][y]
			original, ok := color.RGBAModel.Convert(p).(color.RGBA)

			if ok {
				newImage.Set(x, y, original)
			}
		}
	}

	return newImage

}

func UpsideDownImageTensor(pixels *[][]color.Color)(image.Image){
	p := *pixels
	wg := sync.WaitGroup{}
	rect := image.Rect(0, 0, len(p), len(p[0]))
	newImage := image.NewRGBA(rect)
	for x:= 0; x < len(p); x++{
		for y:= 0; y < len(p[0]); y++{
			if y > (len(p[0]) / 2 + 1) {
				break
			}
			wg.Add(1)
			go func (x, y int){
				pix := p[x][y]
				pix2 := p[len(p)-x-1][len(p[0])-y-1]
				newImage.Set(x, y, pix2)
				newImage.Set(len(p)-x-1, len(p[0])-y-1, pix)
				wg.Done()
			}(x, y)	
		}
	}
	return newImage
}

func RotateImage90AntiClock(img *image.Image) image.Image{

	nrotated := imaging.Rotate90(*img)
	return nrotated
}

func FitImage(img *image.Image, toWidth int, toHeight int) image.Image{
	srcImage := imaging.Fit(*img, toWidth, toHeight, imaging.Lanczos)
	dstImage := image.NewRGBA(image.Rect(0, 0, toWidth, toHeight))

	draw.CatmullRom.Scale(dstImage, dstImage.Rect, srcImage, srcImage.Rect, draw.Over, nil) //convert back to image.RGBA

	return dstImage 
}

func EncodeImageAsJpeg(img image.Image, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	
	err = jpeg.Encode(f, img, nil)
	if err != nil {
		return err
	}
	return nil
	
}

func ConvertGreyScale(pixels *[][]color.Color)(image.Image){

	p := *pixels
	wg := sync.WaitGroup{}
	rect := image.Rect(0,0,len(p), len(p[0]))
	newImage := image.NewRGBA(rect)
	for x:= 0; x < len(p); x++{
		for y:= 0; y < len(p[0]); y++{
			wg.Add(1)
			go func(x, y int){
				pix := p[x][y]
				originalColor, ok:= color.RGBAModel.Convert(pix).(color.RGBA)
				if !ok {
					log.Fatalf("color.color conversion went wrong")
				}
				grey := uint8(float64(originalColor.R)*0.21 + float64(originalColor.G)*0.72 + float64(originalColor.B)*0.07)
				col := color.RGBA{
					grey,
					grey,
					grey,
					originalColor.A,
				}
				newImage.Set(x, y, col)
				
				wg.Done()
			}(x, y)
		} 
	} 
	return newImage
}


func ScaleImage(srcImg *image.Image, toWidth int, toHeight int, scaleQuality int) (image.Image) {
	dstImage := image.NewRGBA(image.Rect(0, 0, toWidth, toHeight))
	
	switch(scaleQuality){
	case ScaleBestQ:
		draw.CatmullRom.Scale(dstImage, dstImage.Rect, *srcImg, (*srcImg).Bounds(), draw.Over, nil)
	case ScaleBetterQ:
		draw.BiLinear.Scale(dstImage, dstImage.Rect, *srcImg, (*srcImg).Bounds(), draw.Over, nil)
	case ScaleGoodQ:
		draw.ApproxBiLinear.Scale(dstImage, dstImage.Rect, *srcImg, (*srcImg).Bounds(), draw.Over, nil)
	default: 
		draw.NearestNeighbor.Scale(dstImage, dstImage.Rect, *srcImg, (*srcImg).Bounds(), draw.Over, nil)
	}
	return dstImage
}


func OrientateAndFitImage(img *image.Image, toWidth int, toHeight int) image.Image {
	size:= (*img).Bounds().Size()
	var result image.Image
	if size.X > size.Y{
		rotated := RotateImage90AntiClock(img) // to become potrait
		result = FitImage(&rotated, toWidth, toHeight)
	} else {
		result = FitImage(img, toWidth, toHeight)
	}
	return result
}
