package imageutil

import (
	"os"
	"testing"
)

func TestOpenImageIntoTensorAndBack(t *testing.T) {
	img, err := OpenImage("./test/test.jpg")
	if err != nil {
		t.Fatal(err)
	}

	pixels := GetImageTensor(img)
	if pixels[0] == nil {
		t.Fatalf("tensor is empty; failed to convert from image")
	}

	outImg := GetBackImage(&pixels)
	
	err = EncodeImageAsJpeg(outImg, "./test/test_out.jpg")
	if err != nil {
		t.Fatal(err)
	}

	f, err := os.Open("./test/test_out.jpg")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	fi, _ := f.Stat()
	realSize := fi.Size() > 0
	if realSize == false {
		t.Errorf("Zero size image obtained back, encode fail")
	}

}

func TestCovertImageIntoGreyscale(t *testing.T){
	img, err := OpenImage("./test/test.jpg")
	if err != nil {
		t.Fatal(err)
	}

	pixels := GetImageTensor(img)
	if pixels[0] == nil {
		t.Fatalf("tensor is empty; failed to convert from image")
	}

	outGrey := ConvertGreyScale(&pixels)
	
	err = EncodeImageAsJpeg(outGrey, "./test/test_grey_out.jpg")
	if err != nil {
		t.Fatal(err)
	}

	f, err := os.Open("./test/test_grey_out.jpg")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	fi, _ := f.Stat()
	realSize := fi.Size() > 0
	if realSize == false {
		t.Errorf("Zero size image obtained back, encode fail")
	}
}

func TestScaleImage(t *testing.T){
	img, err := OpenImage("./test/test.jpg")
	if err != nil {
		t.Fatal(err)
	}

	toWidth := (img.Bounds().Max.X - img.Bounds().Min.X) / 2
	toHeight := (img.Bounds().Max.Y - img.Bounds().Min.Y) / 2
	outScaled := ScaleImage(&img, toWidth, toHeight, ScaleBestQ)
	
	err = EncodeImageAsJpeg(outScaled, "./test/test_scale_halved_out.jpg")
	if err != nil {
		t.Fatal(err)
	}

	f, err := os.Open("./test/test_grey_out.jpg")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	fi, _ := f.Stat()
	realSize := fi.Size() > 0
	if realSize == false {
		t.Errorf("Zero size image obtained back, encode fail")
	}
}

func TestUpsideDownImage(t *testing.T){
	img, err := OpenImage("./test/test.jpg")
	if err != nil {
		t.Fatal(err)
	}


	pixels := GetImageTensor(img)
	if pixels[0] == nil {
		t.Fatalf("tensor is empty; failed to convert from image")
	}

	outUpsideDown := UpsideDownImageTensor(&pixels)
	
	err = EncodeImageAsJpeg(outUpsideDown, "./test/test_upside_down_out.jpg")
	if err != nil {
		t.Fatal(err)
	}

	f, err := os.Open("./test/test_upside_down_out.jpg")
	if err != nil {
		t.Fatal(err)
	}

	fi, _ := f.Stat()
	realSize := fi.Size() > 0
	if realSize == false {
		t.Errorf("Zero size image obtained back, encode fail")
	}

	f.Close()

	//test with image that height is odd in dimension
	oddYImage := ScaleImage(&img, 768, 431, ScaleNormalQ)

	pixels = GetImageTensor(oddYImage)
	if pixels[0] == nil {
		t.Fatalf("tensor is empty; failed to convert from image")
	}

	outUpsideDownOddY := UpsideDownImageTensor(&pixels)
	
	err = EncodeImageAsJpeg(outUpsideDownOddY, "./test/test_upside_down_oddY_out.jpg")
	if err != nil {
		t.Fatal(err)
	}

	f, err = os.Open("./test/test_upside_down_oddY_out.jpg")
	if err != nil {
		t.Fatal(err)
	}


	fi, _ = f.Stat()
	realSize = fi.Size() > 0
	if realSize == false {
		t.Errorf("Zero size image obtained back, encode fail")
	}

	defer f.Close()
}

func TestRotate90Image(t *testing.T){
	img, err := OpenImage("./test/test.jpg")
	if err != nil {
		t.Fatal(err)
	}

	outRotated := RotateImage90AntiClock(&img)
	
	err = EncodeImageAsJpeg(outRotated, "./test/test_rotated_out.jpg")
	if err != nil {
		t.Fatal(err)
	}

	f, err := os.Open("./test/test_rotated_out.jpg")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	fi, _ := f.Stat()
	realSize := fi.Size() > 0
	if realSize == false {
		t.Errorf("Zero size image obtained back, encode fail")
	}
}


func TestFitImage(t *testing.T){
	img, err := OpenImage("./test/test.jpg")
	if err != nil {
		t.Fatal(err)
	}

	outFitted := FitImage(&img, 264, 176)
	
	err = EncodeImageAsJpeg(outFitted, "./test/test_Fitted_out.jpg")
	if err != nil {
		t.Fatal(err)
	}

	f, err := os.Open("./test/test_Fitted_out.jpg")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	fi, _ := f.Stat()
	realSize := fi.Size() > 0
	if realSize == false {
		t.Errorf("Zero size image obtained back, encode fail")
	}
}

func TestOrientateAndFitImageLandscape(t *testing.T){
	img, err := OpenImage("./test/test.jpg")
	if err != nil {
		t.Fatal(err)
	}

	outOrientFitted := OrientateAndFitImage(&img, 176, 264)
	
	err = EncodeImageAsJpeg(outOrientFitted, "./test/test_orient_fit_land_out.jpg")
	if err != nil {
		t.Fatal(err)
	}

	f, err := os.Open("./test/test_orient_fit_land_out.jpg")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	fi, _ := f.Stat()
	realSize := fi.Size() > 0
	if realSize == false {
		t.Errorf("Zero size image obtained back, encode fail")
	}
}

func TestOrientateAndFitImagePortrait(t *testing.T){
	img, err := OpenImage("./test/test_portrait.jpg")
	if err != nil {
		t.Fatal(err)
	}

	outOrientFitted := OrientateAndFitImage(&img, 176, 264)
	
	err = EncodeImageAsJpeg(outOrientFitted, "./test/test_orient_fit_port_out.jpg")
	if err != nil {
		t.Fatal(err)
	}

	f, err := os.Open("./test/test_orient_fit_port_out.jpg")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	fi, _ := f.Stat()
	realSize := fi.Size() > 0
	if realSize == false {
		t.Errorf("Zero size image obtained back, encode fail")
	}
}